# 数据库设置指南

## 错误信息
```
ERROR: permission denied for schema public (SQLSTATE 42501)
```

## 解决方案

### 1. 连接到 PostgreSQL
首先以 postgres 用户身份连接到 PostgreSQL：

```bash
psql -U postgres
```

### 2. 执行以下 SQL 命令

```sql
-- 创建用户
CREATE USER tshirt WITH PASSWORD 'tshirt';

-- 创建数据库
CREATE DATABASE tshirt_db;

-- 授予权限
GRANT ALL PRIVILEGES ON DATABASE tshirt_db TO tshirt;

-- 连接到新创建的数据库
\c tshirt_db

-- 授予 schema 权限
GRANT ALL PRIVILEGES ON SCHEMA public TO tshirt;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO tshirt;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO tshirt;

-- 设置默认权限
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO tshirt;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO tshirt;
```

### 3. 验证连接

```bash
# 测试连接
psql -h localhost -U tshirt -d tshirt_db -c "SELECT 1"
```

### 4. 重新启动应用

完成上述步骤后，重新启动应用程序：

```bash
go run cmd/server/main.go
```

## 配置说明

当前的数据库配置在 `configs/config.yaml` 中：

```yaml
database:
  dsn: "host=localhost user=tshirt password=tshirt dbname=tshirt_db port=5432 sslmode=disable"
```

## 故障排除

1. **确保 PostgreSQL 正在运行**
   ```bash
   sudo systemctl status postgresql
   ```

2. **检查 PostgreSQL 监听地址**
   确保 `postgresql.conf` 中的 `listen_addresses` 包含 `localhost`

3. **检查 pg_hba.conf**
   确保有适当的认证规则：
   ```
   host    all             all             127.0.0.1/32            md5
   host    all             all             ::1/128                 md5
   ```