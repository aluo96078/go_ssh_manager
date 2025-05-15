# Go SSH Manager 使用說明

## 簡介
Go SSH Manager 是一款專為特定任務設計的 SSH 管理工具，採用 Go 語言開發。

## 編譯指南
1. 請確認已安裝 Go 編譯器，版本需為 1.24.1 或以上。
2. 克隆本專案。
3. 打開終端機，切換至專案目錄：
   ```bash
   cd path_to_your_project
   ```
4. 執行以下命令進行編譯並賦予執行權限：
   ```bash
   go build -o <name> ./cmd
   chmod +x <name>
   ```
   編譯完成後，將生成指定名稱的可執行檔。

## 安裝與環境配置

### 從 GitHub Releases 安裝

每當發布新版本時，GitHub Actions 會自動構建並釋出對應平台的可執行檔。請依照以下步驟完成安裝與環境變數配置。

---

### macOS (Intel / Apple Silicon)

#### 下載
```bash
# Intel Mac
curl -L https://github.com/aluo96078/go_ssh_manager/releases/latest/download/gsh-darwin-amd64.zip -o gsh.zip

# Apple Silicon Mac
curl -L https://github.com/aluo96078/go_ssh_manager/releases/latest/download/gsh-darwin-arm64.zip -o gsh.zip
```

#### 解壓縮
```bash
unzip gsh.zip
```

#### 移動並設置執行權限
```bash
chmod +x gsh-darwin-amd64
sudo mv gsh-darwin-amd64 /usr/local/bin/gsh
```

#### 驗證安裝
```bash
gsh
```

---

### Linux

#### 下載
```bash
# x86-64
curl -L https://github.com/aluo96078/go_ssh_manager/releases/latest/download/gsh-linux-amd64.zip -o gsh.zip

# arm
curl -L https://github.com/aluo96078/go_ssh_manager/releases/latest/download/gsh-linux-arm64.zip -o gsh.zip
```

#### 解壓縮
```bash
unzip gsh.zip
```

#### 移動並設置執行權限
```bash
chmod +x gsh-linux-amd64
sudo mv gsh-linux-amd64 /usr/local/bin/gsh
```

#### 驗證安裝
```bash
gsh
```

---

### Windows

#### 下載
請至 [GitHub Releases](https://github.com/aluo96078/go_ssh_manager/releases/latest) 下載 `gsh-windows-amd64.exe.zip`。

#### 解壓縮
使用解壓縮軟體（如 WinRAR 或 7-Zip）解壓縮下載的壓縮檔。

#### 使用方式

##### 方法一：直接執行
1. （可選）將解壓後的 `.exe` 檔案重新命名為 `gsh.exe`。
2. 雙擊執行或在命令提示字元（CMD）中執行。

##### 方法二：添加至系統 PATH
1. 將執行檔移動至固定路徑（例如 `C:\Tools`）。
2. 開啟「系統內容」>「進階系統設定」>「環境變數」。
3. 在「系統變數」中選擇「Path」，點擊「編輯」。
4. 新增執行檔所在目錄（例如 `C:\Tools`）。
5. 確認並關閉所有視窗。

#### 驗證安裝
```cmd
gsh
```

---

### 手動添加至 PATH（適用於自行編譯或其他情況）

若您自行編譯或手動放置可執行檔，請將其所在目錄加入系統 PATH 中：

```bash
export PATH=$PATH:path_to_your_project/<name>
```

為使設定永久生效，請將上述命令加入您的 shell 配置檔（如 `~/.bashrc` 或 `~/.zshrc`）：

```bash
echo 'export PATH=$PATH:path_to_your_project/<name>' >> ~/.bashrc
# 或
echo 'export PATH=$PATH:path_to_your_project/<name>' >> ~/.zshrc
```

## 使用說明

啟動 SSH 管理工具：

```bash
gsh
```

可用命令包括：

- `connect`：連接至已保存的 SSH 設定
- `new`：新增 SSH 設定
- `list`：列出所有已保存的 SSH 設定
- `delete`：刪除指定的 SSH 設定
- `update`：更新既有的 SSH 設定

## 支援

若有任何問題，歡迎聯繫開發者或於 GitHub 提交 Issue。
