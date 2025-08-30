#!/bin/bash

# 自动监听代码变更并重新生成API文档的脚本
# 当检测到Go源文件变更时，自动重新生成Swagger文档

set -e

echo "启动API文档自动生成监听..."
echo "当检测到*.go文件变更时，将自动重新生成文档"
echo "按 Ctrl+C 停止监听"

# 检查是否安装了必要的工具
if ! command -v swag >/dev/null 2>&1; then
    echo "安装 swag 工具..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

if ! command -v inotifywait >/dev/null 2>&1; then
    echo "错误: 需要安装 inotify-tools"
    echo "Ubuntu/Debian: sudo apt-get install inotify-tools"
    echo "CentOS/RHEL: sudo yum install inotify-tools"
    echo "macOS: brew install fswatch (使用 fswatch 替代)"
    exit 1
fi

# 初始生成文档
echo "初始生成API文档..."
swag init -g cmd/server/main.go -o docs/
echo "文档生成完成!"

# 监听文件变更
echo "开始监听Go文件变更..."
while true; do
    # 监听所有Go文件的变更
    inotifywait -r -e modify,create,delete --include='.*\.go$' . 2>/dev/null
    
    echo "检测到Go文件变更，重新生成文档..."
    
    # 添加延迟以避免频繁重生成
    sleep 2
    
    # 重新生成文档
    if swag init -g cmd/server/main.go -o docs/ 2>/dev/null; then
        echo "文档更新完成! $(date)"
    else
        echo "文档生成失败，请检查代码语法"
    fi
    
    echo "继续监听文件变更..."
done