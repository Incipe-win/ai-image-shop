.PHONY: help build run dev test clean docs gen-docs serve-docs

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  build      - 编译应用程序"
	@echo "  run        - 运行应用程序"
	@echo "  dev        - 开发模式运行(使用Air热重载)"
	@echo "  test       - 运行测试"
	@echo "  clean      - 清理构建文件"
	@echo "  docs       - 生成API文档"
	@echo "  gen-docs   - 生成并服务API文档"
	@echo "  serve-docs - 在浏览器中打开API文档"

# 构建应用程序
build:
	@echo "构建应用程序..."
	go build -o bin/server cmd/server/main.go

# 运行应用程序
run: build
	@echo "运行应用程序..."
	./bin/server

# 开发模式运行
dev:
	@echo "开发模式运行..."
	@if [ -f .air.toml ]; then \
		air; \
	else \
		go run cmd/server/main.go; \
	fi

# 运行测试
test:
	@echo "运行测试..."
	go test ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -rf tmp/

# 生成API文档
docs:
	@echo "生成API文档..."
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "安装 swag 工具..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init -g cmd/server/main.go -o docs/

# 生成并打开API文档
gen-docs: docs
	@echo "API文档已生成!"
	@echo "启动服务器后，可以通过以下链接访问API文档:"
	@echo "  Swagger UI: http://localhost:8080/swagger/index.html"
	@echo "  JSON格式:   http://localhost:8080/swagger/doc.json"

# 在浏览器中打开API文档
serve-docs:
	@echo "在浏览器中打开API文档..."
	@if command -v xdg-open >/dev/null 2>&1; then \
		xdg-open http://localhost:8080/swagger/index.html; \
	elif command -v open >/dev/null 2>&1; then \
		open http://localhost:8080/swagger/index.html; \
	else \
		echo "请手动访问: http://localhost:8080/swagger/index.html"; \
	fi

# 安装依赖
deps:
	@echo "下载依赖..."
	go mod download
	go mod tidy

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 检查代码
vet:
	@echo "检查代码..."
	go vet ./...

# 完整的开发流程
dev-setup: deps docs
	@echo "开发环境设置完成!"
	@echo "运行 'make dev' 启动开发服务器"