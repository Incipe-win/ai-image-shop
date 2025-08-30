# 认证系统更新

## 修复的问题

### 1. 注册按钮无法点击问题
- **问题**: 登录页面中的"立即注册"链接无法点击
- **修复**: 在 `static/js/app.js` 中的 `showAuthModal` 方法中添加了事件重新绑定逻辑
- **位置**: `static/js/app.js:104-118`

### 2. JWT Token 失效时间设置
- **问题**: JWT token 过期时间为24小时，太长
- **修复**: 将 token 过期时间设置为1小时
- **位置**: `internal/handler/user_handler.go:78`

### 3. Refresh Token 机制实现
- **新增**: 完整的 refresh token 机制

## 新增功能

### 1. Refresh Token 支持
- **后端**: 添加了 `/api/v1/auth/refresh` 端点
- **数据库**: 在 `users` 表中添加了 `refresh_token` 字段
- **前端**: 自动 token 刷新机制

### 2. Token 自动刷新
- **功能**: 当 access token 过期时，自动使用 refresh token 获取新 token
- **位置**: `static/js/app.js:326-389`

## 代码变更

### 后端 (`internal/handler/user_handler.go`)
1. 添加了 `generateRefreshToken()` 和 `generateTokens()` 辅助函数
2. 修改了 `Register` 和 `Login` 函数以支持 refresh token
3. 新增了 `RefreshToken` 函数处理 token 刷新
4. 更新了响应结构体以包含 refresh token

### 前端 (`static/js/app.js`)
1. 添加了 `refreshToken` 状态管理
2. 实现了 `refreshAuthToken()` 方法
3. 修改了 `apiRequest()` 方法以支持自动 token 刷新
4. 更新了登录/注册成功后的 token 存储逻辑

### 数据库 (`internal/model/user.go`)
1. 添加了 `RefreshToken` 字段到 `User` 结构体

### 路由 (`internal/handler/router.go`)
1. 添加了 `/api/v1/auth/refresh` 路由

## 安全改进

1. **更短的 token 有效期**: 从24小时改为1小时
2. **Refresh token 轮换**: 每次刷新都会生成新的 refresh token
3. **安全的 token 生成**: 使用加密安全的随机数生成 refresh token
4. **数据库存储**: refresh token 存储在数据库中，可撤销

## 使用方式

### 登录/注册响应
```json
{
  "id": 1,
  "username": "user",
  "email": "user@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "abc123def456...",
  "message": "Login successful"
}
```

### Token 刷新请求
```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "abc123def456..."
}
```

### Token 刷新响应
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "new_refresh_token_here...",
  "message": "Tokens refreshed successfully"
}
```

## 测试

启动服务器后，可以测试以下功能：
1. 点击"注册"按钮打开注册模态框
2. 点击"登录"按钮打开登录模态框  
3. 在登录模态框中点击"立即注册"切换到注册
4. 在注册模态框中点击"立即登录"切换到登录
5. 完整的认证流程（注册、登录、token刷新）