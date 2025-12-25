#!/bin/bash

echo "开始测试项目..."

echo "1. 检查Go环境..."
go version

echo ""
echo "2. 下载依赖..."
go mod tidy

echo ""
echo "3. 编译项目..."
go build -o pet-service main.go

if [ $? -eq 0 ]; then
    echo "✓ 编译成功"
    echo ""
    echo "4. 运行服务（测试模式，5秒后自动退出）..."
    timeout 5 ./pet-service || true
else
    echo "✗ 编译失败"
    exit 1
fi

echo ""
echo "测试完成"
