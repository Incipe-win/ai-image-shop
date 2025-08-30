# API 文档说明

本项目使用 Swagger/OpenAPI 2.0 标准自动生成 RESTful API 文档。

## 快速开始

### 1. 生成API文档

```bash
# 使用 Makefile
make docs

# 或直接使用 swag 命令
swag init -g cmd/server/main.go -o docs/
```

### 2. 启动服务器

```bash
# 开发模式
make dev

# 或直接运行
go run cmd/server/main.go
```

### 3. 访问文档

启动服务器后，可以通过以下方式访问API文档：

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON格式**: http://localhost:8080/swagger/doc.json
- **YAML文件**: `docs/swagger.yaml`

## 自动更新功能

### 监听代码变更自动重新生成文档

**Linux 系统:**
```bash
./scripts/docs-watch.sh
```

**macOS 系统:**
```bash
./scripts/docs-watch-mac.sh
```

### 系统要求

**Linux:**
- 需要安装 `inotify-tools`
- Ubuntu/Debian: `sudo apt-get install inotify-tools`
- CentOS/RHEL: `sudo yum install inotify-tools`

**macOS:**
- 需要安装 `fswatch`
- `brew install fswatch`

## API 端点概览

### 认证相关 (Authentication)
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新访问令牌

### 设计相关 (Designs)
- `POST /api/v1/designs/generate` - 生成AI设计 🔒
- `GET /api/v1/designs/my-designs` - 获取用户设计 🔒

### 产品相关 (Products)
- `GET /api/v1/tshirts` - 获取T恤列表

### 系统相关 (System)
- `GET /api/v1/health` - 健康检查

> 🔒 表示需要Bearer token认证

## 认证说明

受保护的API端点需要在请求头中包含Bearer token：

```
Authorization: Bearer <your-jwt-token>
```

## 添加新的API端点文档

### 1. 在处理函数上添加Swagger注解

```go
// ExampleAPI godoc
// @Summary 示例API
// @Description 这是一个示例API的详细描述
// @Tags example
// @Accept json
// @Produce json
// @Param request body ExampleRequest true "请求参数"
// @Success 200 {object} ExampleResponse "成功响应"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Router /example [post]
func ExampleAPI(c *gin.Context) {
    // 处理逻辑
}
```

### 2. 定义请求和响应结构体

```go
type ExampleRequest struct {
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age" binding:"min=0,max=150"`
}

type ExampleResponse struct {
    ID      uint   `json:"id"`
    Message string `json:"message"`
}
```

### 3. 重新生成文档

```bash
make docs
```

## 注解说明

### 常用注解标签

- `@Summary` - API简短描述
- `@Description` - API详细描述
- `@Tags` - API分组标签
- `@Accept` - 接受的内容类型
- `@Produce` - 返回的内容类型
- `@Param` - 参数定义
- `@Success` - 成功响应
- `@Failure` - 错误响应
- `@Router` - 路由定义
- `@Security` - 安全认证要求

### 参数类型

- `query` - URL查询参数
- `path` - URL路径参数
- `header` - 请求头参数
- `body` - 请求体参数
- `formData` - 表单数据

### 示例

```go
// @Param id path int true "用户ID"
// @Param name query string false "用户名称"
// @Param request body UserRequest true "用户信息"
```

## 文档结构

```
docs/
├── docs.go       # 生成的Go代码
├── swagger.json  # JSON格式文档
└── swagger.yaml  # YAML格式文档
```

## 故障排除

### 常见问题

1. **文档生成失败**
   - 检查Go代码语法是否正确
   - 确认Swagger注解格式是否正确
   - 查看错误日志信息

2. **文档内容不更新**
   - 重新运行 `make docs` 命令
   - 重启服务器
   - 清除浏览器缓存

3. **无法访问Swagger UI**
   - 确认服务器正在运行
   - 检查端口是否被占用
   - 验证路由配置是否正确

### 调试技巧

- 使用 `swag init -g cmd/server/main.go -o docs/ --parseVendor` 解析vendor包
- 添加 `--parseDependency` 参数解析依赖包
- 使用 `-d .` 指定搜索目录

## 扩展功能

### 自定义文档主题

可以通过修改Swagger UI配置来自定义文档外观：

```go
// 在router.go中自定义Swagger配置
config := &ginSwagger.Config{
    URL: "doc.json",
    DocExpansion: "list",
    DeepLinking:  true,
}
r.GET("/swagger/*any", ginSwagger.CustomWrapHandler(config, swaggerFiles.Handler))
```

### 多环境配置

可以为不同环境设置不同的文档配置：

```go
// 在main.go中根据环境设置不同的host
if env == "production" {
    docs.SwaggerInfo.Host = "api.yourdomain.com"
} else {
    docs.SwaggerInfo.Host = "localhost:8080"
}
```

## 最佳实践

1. **保持注解更新** - 每次修改API时同步更新Swagger注解
2. **使用有意义的标签** - 合理分组API端点
3. **提供详细描述** - 包含足够的信息帮助API使用者
4. **定义完整的响应结构** - 包括成功和错误响应
5. **使用示例数据** - 在结构体中添加示例值
6. **定期验证文档** - 确保文档与实际API行为一致