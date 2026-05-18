# Air 热重载配置说明

## 问题解决

如果你在使用 `air` 时遇到 PowerShell 执行错误，可以使用以下方法：

### 方法1：使用 CMD 运行（推荐）

使用 `start-dev.bat` 脚本：

```cmd
start-dev.bat
```

或者在 CMD 中直接运行：

```cmd
cd softeng-platform\softeng-platform
air
```

### 方法2：在 PowerShell 中使用 CMD

```powershell
cd softeng-platform\softeng-platform
cmd /c air
```

### 方法3：如果 air 仍然有问题，使用 nodemon 替代

如果 air 在 Windows 下仍然有问题，可以使用 `nodemon`：

```bash
npm install -g nodemon
nodemon --exec "go run cmd/server/main.go" --watch . --ext go
```

### 方法4：使用 Go 的官方工具（最简单）

创建一个简单的监听脚本，或者直接使用：

```bash
# 安装工具
go install github.com/cosmtrek/air@latest

# 确保 air 在 PATH 中，然后运行
air
```

## 配置文件

配置文件位于 `.air.toml`，主要设置：
- 监听 `cmd/`、`internal/`、`pkg/` 目录
- 编译输出到 `tmp/main.exe`
- 自动重新编译和重启

## 注意事项

- Windows 下编译的文件需要 `.exe` 扩展名
- 如果使用 PowerShell，可能需要使用 CMD 模式运行 air
- `tmp/` 目录存储编译后的二进制文件，可以添加到 `.gitignore`

