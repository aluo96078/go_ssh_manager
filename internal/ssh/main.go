package ssh

import (
	"fmt"
	"os"
	"os/exec"
)

type SSHConfig struct {
	ConfigName   string
	HostName     string
	UserName     string
	Port         string
	IdentityFile string
	Comment      string
}

var funcMap = map[string]func() error{
	"connect": ConnectSSH,
	"new":     NewSSHConfig,
	"list":    ListSSHConfig,
	"delete":  DeleteSSHConfig,
	"update":  UpdateSSHConfig,
}

func ProcessSSHConfig(action string) (string, error) {
	if fn, ok := funcMap[action]; ok {
		if err := fn(); err != nil {
			return "", err
		}
	}
	return "", nil
}

/*
新建 SSH 設定
這個函數會讓用戶輸入 SSH 設定的各個參數，然後將這些參數寫入 SSH 設定文件中。
用戶可以選擇是否使用私鑰文件進行驗證，並且可以選擇是否輸入備註。
這個函數的主要步驟如下：
 1. 輸入設定檔名稱
 2. 輸入用戶名
 3. 輸入主機名
 4. 輸入端口號
 5. 輸入私鑰文件路徑
 6. 輸入私鑰文件密碼
 7. 輸入備註
 8. 輸出指令
 9. 輸出設定檔名稱、用戶名、主機名、端口號、密鑰文件、備註
 10. 輸出提示信息
 11. 等待用戶確認
 12. 獲取 SSH 設定列表
 13. 將新的 SSH 設定添加到列表中
 14. 寫入 SSH 設定文件
 15. 輸出寫入成功的提示信息
 16. 返回 nil，表示新建成功
*/
func NewSSHConfig() error {

	// 輸入設定檔名稱
	var configName string
	err := inputVariable(&configName, "請輸入設定檔名稱: ", false, "")
	if err != nil {
		return fmt.Errorf("輸入設定檔名稱失敗: %v", err)
	}
	// 判斷設定檔名稱是否重複
	if _, err := getSSHConfigByName(configName); err == nil {
		return fmt.Errorf("設定檔名稱 %s 已經存在", configName)
	}

	// 輸入用戶名
	var username string = "root"
	err = inputVariable(&username, "請輸入用戶名(默認 root): ", true, "root")
	if err != nil {
		return fmt.Errorf("輸入用戶名失敗: %v", err)
	}

	// 輸入主機名
	var hostname string
	err = inputVariable(&hostname, "請輸入主機名: ", false, "")
	if err != nil {
		return fmt.Errorf("輸入主機名失敗: %v", err)
	}

	// 輸入端口號
	var port string = "22"
	err = inputVariable(&port, "請輸入端口號(默認 22): ", true, "22")
	if err != nil {
		return fmt.Errorf("輸入端口號失敗: %v", err)
	}

	// 非對稱加密設定
	var useSSHKeygen string
	var keyPath string
	var keyPass string
	fmt.Print("是否使用私鑰文件(非對稱加密)進行驗證 (y/n): ")
	fmt.Scanln(&useSSHKeygen)
	if useSSHKeygen == "y" {
		keyPath, keyPass, err = selectSSHKeyFile()
		if err != nil {
			return err
		}
	}

	// 輸入備註
	var comment string
	err = inputVariable(&comment, "請輸入備註: ", true, "")
	if err != nil {
		return fmt.Errorf("輸入備註失敗: %v", err)
	}

	// 輸出指令
	fmt.Println("以下是生成的 SSH 設定:")
	fmt.Println("=======================================")
	fmt.Printf("ssh %s@%s -p %s -i %s\n", username, hostname, port, keyPath)
	fmt.Println("=======================================")
	fmt.Printf("設定檔名稱: %s\n", configName)
	fmt.Printf("用戶名: %s\n", username)
	fmt.Printf("主機名: %s\n", hostname)
	fmt.Printf("端口號: %s\n", port)
	fmt.Printf("密鑰文件: %s\n", keyPath)
	fmt.Printf("密鑰文件密碼: %s\n", keyPass)
	fmt.Printf("備註: %s\n", comment)
	fmt.Println("=======================================")
	fmt.Println("請確認以上信息是否正確，然後按 Enter 鍵保存，或按 Ctrl+C 退出。")
	fmt.Scanln()

	sshConfigs, err := getSSHConfigList()
	if err != nil {
		return err
	}

	newConfig := SSHConfig{
		ConfigName:   configName,
		HostName:     hostname,
		UserName:     username,
		Port:         port,
		IdentityFile: keyPath,
		Comment:      comment,
	}

	sshConfigs = append(sshConfigs, newConfig)

	err = rewriteSSHConfigFile(sshConfigs)
	if err != nil {
		return fmt.Errorf("寫入 SSH 設定文件失敗: %v", err)
	}

	fmt.Println("SSH 設定文件寫入成功")
	return nil
}

/*
列出 SSH 設定
這個函數會從 SSH 設定文件中獲取所有的 SSH 設定，然後將它們列出來。
用戶可以選擇要查看的 SSH 設定，然後函數會將選擇的 SSH 設定輸出到終端。
這個函數的主要步驟如下：
 1. 獲取 SSH 設定列表
 2. 如果沒有找到任何 SSH 設定，則輸出提示信息
 3. 列出所有的 SSH 設定
 4. 返回 nil，表示列出成功
*/
func ListSSHConfig() error {
	sshConfigs, err := getSSHConfigList()
	if err != nil {
		return err
	}
	if len(sshConfigs) == 0 {
		fmt.Println("沒有找到任何 SSH 設定")
		return nil
	}
	// 輸出 SSH 設定
	fmt.Println("以下是 SSH 設定:")
	fmt.Println("=======================================")
	for _, config := range sshConfigs {
		fmt.Printf("Host: %s\n", config.ConfigName)
		fmt.Printf("  HostName: %s\n", config.HostName)
		fmt.Printf("  User: %s\n", config.UserName)
		fmt.Printf("  Port: %s\n", config.Port)
		fmt.Printf("  IdentityFile: %s\n", config.IdentityFile)
		fmt.Printf("  Comment: %s\n", config.Comment)
		fmt.Println("=======================================")
	}
	return nil
}

/*
刪除 SSH 設定
這個函數會從 SSH 設定文件中獲取所有的 SSH 設定，然後讓用戶選擇要刪除的設定。
用戶可以選擇要刪除的設定，然後函數會將刪除後的 SSH 設定寫入 SSH 設定文件。
這個函數的主要步驟如下：
 1. 獲取 SSH 設定列表
 2. 如果沒有找到任何 SSH 設定，則輸出提示信息
 3. 讓用戶選擇要刪除的 SSH 設定
 4. 獲取選擇的 SSH 設定
 5. 將選擇的 SSH 設定從 SSH 設定列表中刪除
 6. 將更新後的 SSH 設定寫入 SSH 設定文件
 7. 輸出刪除成功的提示信息
 8. 返回 nil，表示刪除成功
*/
func DeleteSSHConfig() error {
	sshConfigs, err := getSSHConfigList()
	if err != nil {
		return err
	}
	if len(sshConfigs) == 0 {
		fmt.Println("沒有找到任何 SSH 設定")
	}
	target, err := selectConfigName(sshConfigs)
	if err != nil {
		return err
	}
	latest := []SSHConfig{}
	for _, config := range sshConfigs {
		if config.ConfigName != target {
			latest = append(latest, config)
		}
	}

	err = rewriteSSHConfigFile(latest)
	if err != nil {
		return fmt.Errorf("寫入 SSH 設定文件失敗: %v", err)
	}

	fmt.Printf("SSH 設定 %s 刪除成功\n", target)
	fmt.Println("=======================================")
	return nil
}

/*
更新 SSH 設定
這個函數會從 SSH 設定文件中獲取所有的 SSH 設定，然後讓用戶選擇要更新的設定。
用戶可以選擇要更新的設定，然後輸入新的設定值。最後，函數會將更新後的 SSH 設定寫入 SSH 設定文件。
這個函數的主要步驟如下：
 1. 獲取 SSH 設定列表
 2. 如果沒有找到任何 SSH 設定，則輸出提示信息
 3. 讓用戶選擇要更新的 SSH 設定
 4. 獲取選擇的 SSH 設定
 5. 詢問用戶是否要變更設定檔名稱、用戶名、主機名、端口號、密鑰文件和備註
 6. 如果用戶選擇變更密鑰文件，則讓用戶選擇新的密鑰文件
 7. 修改選擇的 SSH 設定
 8. 將更新後的 SSH 設定寫入 SSH 設定文件
 9. 輸出更新成功的提示信息
 10. 返回 nil，表示更新成功
*/
func UpdateSSHConfig() error {
	sshConfigs, err := getSSHConfigList()
	if err != nil {
		return err
	}
	if len(sshConfigs) == 0 {
		fmt.Println("沒有找到任何 SSH 設定")
	}
	target, err := selectConfigName(sshConfigs)
	if err != nil {
		return err
	}
	configData, err := getSSHConfigByName(target)
	if err != nil {
		return err
	}
	// 詢問變更設定檔名稱
	err = inputVariable(&configData.ConfigName, "請輸入欲變更的設定檔名稱，爲空則不變更: ", false, configData.ConfigName)
	if err != nil {
		return fmt.Errorf("輸入設定檔名稱失敗: %v", err)
	}
	// 詢問變更用戶名
	err = inputVariable(&configData.UserName, "請輸入欲變更的用戶名，爲空則不變更: ", false, configData.UserName)
	if err != nil {
		return fmt.Errorf("輸入用戶名失敗: %v", err)
	}
	// 詢問變更主機名
	err = inputVariable(&configData.HostName, "請輸入欲變更的主機名，爲空則不變更: ", false, configData.HostName)
	if err != nil {
		return fmt.Errorf("輸入主機名失敗: %v", err)
	}
	// 詢問變更端口號
	err = inputVariable(&configData.Port, "請輸入欲變更的端口號，爲空則不變更: ", false, configData.Port)
	if err != nil {
		return fmt.Errorf("輸入端口號失敗: %v", err)
	}
	// 詢問變更密鑰文件
	var useSSHKeygen string
	fmt.Print("是否變更密鑰文件 (y/n): ")
	fmt.Scanln(&useSSHKeygen)
	if useSSHKeygen == "y" {
		keyPath, _, err := selectSSHKeyFile()
		if err != nil {
			return err
		}
		configData.IdentityFile = keyPath
	}
	// 詢問變更備註
	err = inputVariable(&configData.Comment, "請輸入欲變更的備註，爲空則不變更: ", false, configData.Comment)
	if err != nil {
		return fmt.Errorf("輸入備註失敗: %v", err)
	}
	// 修改設定檔
	latest := []SSHConfig{}
	for _, config := range sshConfigs {
		if config.ConfigName == target {
			config.ConfigName = configData.ConfigName
			config.HostName = configData.HostName
			config.UserName = configData.UserName
			config.Port = configData.Port
			config.IdentityFile = configData.IdentityFile
			config.Comment = configData.Comment
		}
		latest = append(latest, config)
	}
	// 寫入 SSH 設定文件
	err = rewriteSSHConfigFile(latest)
	if err != nil {
		return fmt.Errorf("寫入 SSH 設定文件失敗: %v", err)
	}
	fmt.Printf("SSH 設定 %s 更新成功\n", target)
	fmt.Println("=======================================")
	return nil
}

/*
連接 SSH
這個函數會從 SSH 設定文件中獲取所有的 SSH 設定，然後讓用戶選擇要連接的設定。
用戶可以選擇要連接的設定，然後函數會將選擇的 SSH 設定輸出到終端。
這個函數的主要步驟如下：
 1. 獲取 SSH 設定列表
 2. 如果沒有找到任何 SSH 設定，則輸出提示信息
 3. 讓用戶選擇要連接的 SSH 設定
 4. 獲取選擇的 SSH 設定
 5. 輸出連接指令
 6. 執行 SSH 命令
 7. 如果執行成功，則返回 nil，表示連接成功
 8. 如果執行失敗，則返回錯誤信息
*/
func ConnectSSH() error {
	sshConfigs, err := getSSHConfigList()
	if err != nil {
		return err
	}
	if len(sshConfigs) == 0 {
		fmt.Println("沒有找到任何 SSH 設定")
	}
	target, err := selectConfigName(sshConfigs)
	if err != nil {
		return err
	}
	configData, err := getSSHConfigByName(target)
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("ssh %s@%s -p %s -i %s\n", configData.UserName, configData.HostName, configData.Port, configData.IdentityFile)
	// 執行 SSH 命令
	cmdExec := exec.Command("sh", "-c", cmd)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	err = cmdExec.Run()
	if err != nil {
		return fmt.Errorf("執行 SSH 命令失敗: %v", err)
	}

	return nil
}
