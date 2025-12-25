@echo off
REM 生成GORM代码
echo 正在生成GORM代码...
go run script/gen.go

REM 运行服务
echo 启动服务...
go run main.go
