# PowerShell脚本：使用air启动开发服务器
# 使用方法：.\start-dev.ps1

Write-Host "Starting Go server with hot reload..." -ForegroundColor Green
Write-Host "Directory: $(Get-Location)" -ForegroundColor Gray
Write-Host ""

# 进入脚本所在目录
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptPath

# 检查 air 是否安装（推荐 air-verse 维护版）
$airPath = Get-Command air -ErrorAction SilentlyContinue
if (-not $airPath) {
    Write-Host "Error: 'air' command not found. Please install:" -ForegroundColor Red
    Write-Host "  go install github.com/air-verse/air@latest" -ForegroundColor Yellow
    Write-Host "  并将 %USERPROFILE%\go\bin 加入 PATH" -ForegroundColor Yellow
    exit 1
}

# Windows 下用 cmd 调 air，避免部分 PowerShell 与 fsnotify 兼容问题（见 README-AIR.md）
Write-Host "Running air (via cmd)..." -ForegroundColor Cyan
cmd /c air

