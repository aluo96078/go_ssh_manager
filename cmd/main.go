package main

import (
	"fmt"
	"go_ssh_manager/config"
	"go_ssh_manager/internal/ssh"

	"github.com/manifoldco/promptui"
)

func main() {

	items := make([]string, 0, len(config.Commands))
	for key := range config.Commands {
		items = append(items, key)
	}

	prompt := promptui.Select{
		Label: "請選擇行爲",
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
		return
	}

	_, err = ssh.ProcessSSHConfig(config.Commands[result])
	if err != nil {
		fmt.Printf("執行命令時出錯: \n - %v", err)
		return
	}
}
