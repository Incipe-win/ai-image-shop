#!/bin/bash

# 开发环境启动脚本
# 使用 Air 进行热重载开发

echo "🚀 启动 AI 创意设计工坊开发服务器..."
echo "📁 工作目录: $(pwd)"
echo "🔄 使用 Air 进行热重载开发"
echo ""

# 检查 Air 是否安装
if ! command -v air &> /dev/null; then
    echo "❌ Air 未安装，请先安装: go install github.com/cosmtrek/air@latest"
    exit 1
fi

# 检查配置文件
if [ ! -f ".air.toml" ]; then
    echo "❌ 找不到 .air.toml 配置文件"
    exit 1
fi

echo "✅ Air 已安装"
echo "✅ 配置文件就绪"
echo "🔧 启动热重载服务器..."
echo ""

# 启动 Air
air