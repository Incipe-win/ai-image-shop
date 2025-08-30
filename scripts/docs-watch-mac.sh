#!/bin/bash

# macOS版本的自动监听代码变更并重新生成API文档的脚本
# 使用fswatch替代inotifywait

set -e

echo "启动API文档自动生成监听 (macOS版本)..."
echo "当检测到*.go文件变更时，将自动重新生成文档"
echo "按 Ctrl+C 停止监听"

# 检查是否安装了必要的工具
if ! command -v swag >/dev/null 2>&1; then
    echo "安装 swag 工具..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

if ! command -v fswatch >/dev/null 2>&1; then
    echo "错误: 需要安装 fswatch"
    echo "macOS: brew install fswatch"
    exit 1
fi

# 初始生成文档
echo "初始生成API文档..."
swag init -g cmd/server/main.go -o docs/
echo "文档生成完成!"

# 监听文件变更
echo "开始监听Go文件变更..."
fswatch -o . --include='.*\.go$' | while read num; do
    echo "检测到Go文件变更，重新生成文档..."
    
    # 添加延迟以避免频繁重生成
    sleep 2
    
    # 重新生成文档
    if swag init -g cmd/server/main.go -o docs/ 2>/dev/null; then
        echo "文档更新完成! $(date)"
    else
        echo "文档生成失败，请检查代码语法"
    fi
done