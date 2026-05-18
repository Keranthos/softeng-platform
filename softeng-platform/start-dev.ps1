# PowerShell脚本：使用air启动开发服务器
# 使用方法：.\start-dev.ps1

Write-Host "Starting Go server with hot reload..." -ForegroundColor Green
Write-Host "Directory: $(Get-Location)" -ForegroundColor Gray
Write-Host ""

# 进入脚本所在目录
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptPath

# 检查air是否安装
$airPath = Get-Command air -ErrorAction SilentlyContinue
if (-not $airPath) {
    Write-Host "Error: 'air' command not found. Please install it first:" -ForegroundColor Red
    Write-Host "  go install github.com/cosmtrek/air@latest" -ForegroundColor Yellow
    exit 1
}

# 运行air
Write-Host "Running air..." -ForegroundColor Cyan
air

