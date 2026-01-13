# 热重载配置说明

## 问题说明

`air` 工具在 Windows PowerShell 下可能存在兼容性问题，建议**使用 CMD 而不是 PowerShell** 来运行 `air`。

## 推荐使用方法

### 方法1：使用批处理脚本（最简单）✅

直接双击运行：
```
start-dev.bat
```

或在 CMD 中：
```cmd
cd softeng-platform\softeng-platform
start-dev.bat
```

### 方法2：在 CMD 中直接运行

```cmd
cd softeng-platform\softeng-platform
air
```

### 方法3：在 PowerShell 中使用 CMD 模式

```powershell
cd softeng-platform\softeng-platform
cmd /c air
```

## 如果 air 未安装

```bash
go install github.com/cosmtrek/air@latest
```

确保 `$GOPATH/bin` 或 `$HOME/go/bin` 在系统 PATH 中。

## 配置文件

配置文件：`.air.toml`

- 监听目录：`cmd/`、`internal/`、`pkg/`
- 编译输出：`tmp/main.exe`
- 自动重新编译和重启

## 工作原理

1. `air` 监听 `.go` 文件的变化
2. 文件修改后，自动执行 `go build`
3. 编译成功后，自动重启服务器
4. 编译失败时，显示错误信息，不重启

## 注意事项

- ⚠️ **强烈建议在 CMD 中运行，而不是 PowerShell**
- 修改代码后，等待 1-2 秒让 air 检测到变化
- 按 `Ctrl+C` 停止 air
- `tmp/` 目录会存储编译后的二进制文件

