@echo off
REM Windows批处理脚本：使用air启动开发服务器（CMD模式，更兼容）
REM 使用方法：直接双击此文件，或在CMD中运行：start-dev.bat

echo ========================================
echo   Starting Go Server with Hot Reload
echo   Using: air
echo ========================================
echo.

REM 进入脚本所在目录
cd /d "%~dp0"

REM 检查air是否可用
where air >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] 'air' command not found!
    echo.
    echo Please install air first:
    echo   go install github.com/cosmtrek/air@latest
    echo.
    echo Make sure Go bin directory is in your PATH.
    echo.
    pause
    exit /b 1
)

REM 显示air版本
echo [INFO] Air version:
air -v
echo.

REM 创建tmp目录（如果不存在）
if not exist "tmp" mkdir tmp

REM 运行air（在CMD下运行，避免PowerShell兼容性问题）
echo [INFO] Starting air...
echo [INFO] Air will watch for file changes and auto-reload.
echo [INFO] Press Ctrl+C to stop.
echo.
air

