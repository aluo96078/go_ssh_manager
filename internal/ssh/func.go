package ssh

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func selectSSHKeyFile() (string, string, error) {
	var keyPath, keyPass string
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("獲取用戶主目錄失敗: %v", err)
	}
	files, err := os.ReadDir(home + "/.ssh")
	if err != nil {
		return "", "", fmt.Errorf("讀取目錄失敗: %v", err)
	}
	fileNames := []string{"自行輸入路徑"}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == "known_hosts" || file.Name() == "config" || file.Name() == "authorized_keys" {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}
	prompt := promptui.Select{
		Label: "請選擇密鑰文件",
		Items: fileNames,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("選擇密鑰文件時出錯: %v", err)
	}
	if result == "自行輸入路徑" {
		err = inputVariable(&keyPath, "請輸入密鑰文件路徑: ", false, "")
		if err != nil {
			return "", "", fmt.Errorf("輸入密鑰文件路徑失敗: %v", err)
		}
	} else {
		keyPath = home + "/.ssh/" + result
	}
	if _, err := os.Stat(keyPath); err != nil {
		return "", "", fmt.Errorf("密鑰文件不存在: %s", keyPath)
	}
	err = inputVariable(&keyPass, "請輸入密鑰文件密碼: ", true, "")
	if err != nil {
		return "", "", fmt.Errorf("輸入密鑰文件密碼失敗: %v", err)
	}
	return keyPath, keyPass, nil
}

func rewriteSSHConfigFile(configs []SSHConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("獲取用戶主目錄失敗: %v", err)
	}
	sshConfigPath := home + "/.ssh/config"
	file, err := os.OpenFile(sshConfigPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("打開 SSH 設定文件失敗: %v", err)
	}
	defer file.Close()
	for _, config := range configs {
		sshConfig := fmt.Sprintf("Host %s\n", config.ConfigName)
		sshConfig += fmt.Sprintf("  HostName %s\n", config.HostName)
		sshConfig += fmt.Sprintf("  User %s\n", config.UserName)
		sshConfig += fmt.Sprintf("  Port %s\n", config.Port)
		sshConfig += fmt.Sprintf("  IdentityFile %s\n", config.IdentityFile)
		sshConfig += "  IdentitiesOnly yes\n"
		sshConfig += fmt.Sprintf("  # %s\n", config.Comment)
		_, err = file.WriteString(sshConfig)
		if err != nil {
			return fmt.Errorf("寫入 SSH 設定文件失敗: %v", err)
		}
	}
	return nil
}

func getSSHConfigList() ([]SSHConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("獲取用戶主目錄失敗: %v", err)
	}
	sshConfigPath := home + "/.ssh/config"
	file, err := os.Open(sshConfigPath)
	if err != nil {
		return nil, fmt.Errorf("打開 SSH 設定文件失敗: %v", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("獲取 SSH 設定文件信息失敗: %v", err)
	}
	if fileInfo.Size() == 0 {
		return []SSHConfig{}, nil
	}
	configFile, err := os.ReadFile(sshConfigPath)
	if err != nil {
		return nil, fmt.Errorf("讀取 SSH 設定文件失敗: %v", err)
	}
	lines := strings.Split(string(configFile), "\n")
	var sshConfigs []SSHConfig
	var currentHost string
	current := SSHConfig{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "Host ") {
			if currentHost != "" {
				sshConfigs = append(sshConfigs, current)
			}
			currentHost = strings.TrimSpace(strings.TrimPrefix(trimmed, "Host"))
			current = SSHConfig{}
			current.ConfigName = currentHost
		} else if strings.HasPrefix(trimmed, "HostName ") {
			current.HostName = strings.TrimSpace(strings.TrimPrefix(trimmed, "HostName"))
		} else if strings.HasPrefix(trimmed, "User ") {
			current.UserName = strings.TrimSpace(strings.TrimPrefix(trimmed, "User"))
		} else if strings.HasPrefix(trimmed, "Port ") {
			current.Port = strings.TrimSpace(strings.TrimPrefix(trimmed, "Port"))
		} else if strings.HasPrefix(trimmed, "IdentityFile ") {
			current.IdentityFile = strings.TrimSpace(strings.TrimPrefix(trimmed, "IdentityFile"))
		} else if strings.HasPrefix(trimmed, "# ") {
			current.Comment += strings.TrimSpace(strings.TrimPrefix(trimmed, "#")) + "\n"
		}
	}
	if currentHost != "" {
		sshConfigs = append(sshConfigs, current)
	}

	return sshConfigs, nil
}

func getSSHConfigByName(name string) (*SSHConfig, error) {
	configs, err := getSSHConfigList()
	if err != nil {
		return nil, fmt.Errorf("獲取 SSH 設定列表失敗: %v", err)
	}
	for _, config := range configs {
		if config.ConfigName == name {
			return &config, nil
		}
	}
	return nil, fmt.Errorf("找不到名爲 %s 的 SSH 設定", name)
}

func selectConfigName(configs []SSHConfig) (string, error) {
	items := make([]string, 0, len(configs))
	for _, config := range configs {
		items = append(items, config.ConfigName)
	}
	prompt := promptui.Select{
		Label: "請選擇 SSH 設定",
		Items: items,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("選擇 SSH 設定時出錯: %v", err)
	}
	return result, nil
}

func inputVariable(variable *string, prompt string, empty bool, defaultValue string) error {
	fmt.Print(prompt)
	fmt.Scanln(variable)
	if *variable == "" && empty {
		*variable = defaultValue
	}
	if *variable == "" && !empty {
		return fmt.Errorf("%s 不能為空", prompt)
	}
	return nil
}
