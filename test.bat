@echo off
echo 开始测试项目...

echo 1. 检查Go环境...
go version

echo.
echo 2. 下载依赖...
go mod tidy

echo.
echo 3. 编译项目...
go build -o pet-service.exe main.go

if %errorlevel% == 0 (
    echo 编译成功
    echo.
    echo 4. 运行服务（测试模式，5秒后自动退出）...
    start /b timeout /t 5
    .\pet-service.exe
) else (
    echo 编译失败
    exit /b 1
)

echo.
echo 测试完成
