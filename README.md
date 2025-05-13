# Go SSH Manager 使用文檔

## 簡介
Go SSH Manager 是一個用於處理特定任務的工具，使用 Go 語言編寫。

## 編譯指南
1. 確保已安裝 Go 編譯器，並且版本為 1.24.1 或更高。
2. 克隆此專案
3. 打開終端，導航到項目目錄：
   ```bash
   cd path_to_your_project
   ```
4. 執行以下命令以編譯項目並給予執行權限：
   ```bash
   go build -o <name> ./cmd
   chmod +x <name>
   ```
   編譯完成後，將生成一個指定名稱可執行文件。

## 添加環境參數
1. 將生成的可執行文件添加到系統的 PATH 中：
   ```bash
   export PATH=$PATH:path_to_your_project/<name>
   ```
2. 為了使該設置永久生效，可以將上述命令添加到 `~/.bashrc` 或 `~/.zshrc` 文件中：
   ```bash
   echo 'export PATH=$PATH:path_to_your_project/<name>' >> ~/.bashrc
   # 或者
   echo 'export PATH=$PATH:path_to_your_project/<name>' >> ~/.zshrc
   ```

## 使用方法
編譯並設置環境參數後，可以在終端中直接運行：
```bash
<name>
```

## 支援
如有任何問題，請聯繫開發者或提交 Issue。
