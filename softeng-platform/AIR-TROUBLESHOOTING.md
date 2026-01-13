# Air 热重载问题排查

## 问题：Air 显示 "Process Exit with Code: 1" 但服务器实际在运行

### 现象
- Air 启动程序后，显示 "Process Exit with Code: 1"
- 但服务器实际上在运行（端口8080被监听）
- API 请求可以正常响应

### 原因
这是 Windows 下 air 的已知兼容性问题：
1. Air 在检测长时间运行的进程时，可能会误报退出代码
2. 程序在 goroutine 中启动服务器，主 goroutine 等待信号，air 可能误判为程序已退出
3. 服务器实际上在正常运行，只是 air 的检测机制有问题

### 验证服务器是否真的在运行

**方法1：检查端口**
```powershell
netstat -ano | findstr :8080
```
如果看到 LISTENING 状态，说明服务器在运行。

**方法2：测试 API**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/course?limit=5" -Method GET
```
如果返回 200 状态码，说明服务器正常。

### 解决方案

#### 方案1：忽略 air 的误报（推荐）
如果服务器实际上在运行并响应请求，可以忽略 air 的退出代码显示。
- 服务器在正常监听
- API 可以正常访问
- 代码修改后 air 会自动重新编译和重启

#### 方案2：使用 CMD 而不是 PowerShell
在 CMD 中运行 air，而不是 PowerShell：
```cmd
cd C:\Users\wanweijie\Desktop\shixun\softeng-platform\softeng-platform
air
```
或者使用批处理脚本：
```cmd
start-dev.bat
```

#### 方案3：直接运行程序（不使用 air）
如果 air 一直有问题，可以直接运行程序：
```powershell
cd C:\Users\wanweijie\Desktop\shixun\softeng-platform\softeng-platform
go run cmd/server/main.go
```
注意：这种方式不会自动重新编译，需要手动重启。

### 当前状态

✅ 服务器代码已优化：
- 添加了错误检查
- 改进了信号处理
- 添加了启动确认日志

✅ Air 配置已优化：
- 移除了可能导致问题的 `full_bin` 配置
- 调整了延迟时间

### 建议

如果 air 显示退出代码 1，但服务器实际上在运行：
1. **先验证**：检查端口是否被监听，测试 API 是否响应
2. **如果正常**：可以继续使用 air，忽略退出代码显示
3. **如果不正常**：查看 `tmp/errors.log` 文件，检查是否有真正的错误

### 其他问题

如果遇到其他问题（如端口占用、编译错误等），请检查：
- `tmp/errors.log` - air 的构建日志
- 终端输出 - 程序的运行日志
- 端口占用：`netstat -ano | findstr :8080`

