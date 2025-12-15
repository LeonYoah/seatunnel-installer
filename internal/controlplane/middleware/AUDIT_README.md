# 审计日志中间件

## 概述

审计日志中间件用于记录所有写操作（POST、PUT、PATCH、DELETE）的详细信息，确保系统操作的可追溯性和安全性。

## 功能特性

### 1. 自动记录写操作
- 只记录写操作（POST、PUT、PATCH、DELETE）
- 读操作（GET、HEAD、OPTIONS）不会被记录
- 支持异步写入，不影响主请求性能

### 2. 完整的操作信息
记录的信息包括：
- 操作者（用户ID）
- 租户ID（多租户隔离）
- 操作类型（create/update/delete）
- 资源类型（从URL路径提取）
- 资源ID（如果可以从URL提取）
- 操作详情（请求方法、路径、查询参数、用户代理、客户端IP）
- 操作结果（success/failure）
- 错误信息（如果操作失败）
- 操作时间

### 3. 安全性考虑
- 自动检测敏感数据（password、secret、token、key、credential）
- 敏感数据不会被记录到审计日志中
- 支持租户隔离，确保数据安全

### 4. 智能资源识别
- 自动从URL路径提取资源类型（如 /api/v1/hosts -> hosts）
- 自动识别UUID格式的资源ID
- 支持嵌套路径的资源提取

## 使用方法

### 1. 在路由中添加中间件

```go
// 在所有v1 API路由中添加审计中间件
v1 := router.Group("/api/v1")
v1.Use(middleware.AuditMiddleware(repoManager, logger))
```

### 2. 确保认证中间件在前
审计中间件依赖认证中间件设置的用户信息：

```go
v1.Use(middleware.AuthMiddleware(jwtService, repoManager))  // 必须在前
v1.Use(middleware.AuditMiddleware(repoManager, logger))     // 在后
```

### 3. 查询审计日志

通过审计日志API查询：

```bash
# 获取租户的所有审计日志
GET /api/v1/audit/logs

# 获取指定用户的审计日志
GET /api/v1/audit/users/{user_id}/logs

# 获取指定资源的审计日志
GET /api/v1/audit/resources/{resource}/logs
```

## API接口

### 获取审计日志列表
- **URL**: `GET /api/v1/audit/logs`
- **参数**:
  - `page`: 页码（默认1）
  - `size`: 每页大小（默认20，最大100）
  - `user_id`: 用户ID筛选（可选）
  - `resource`: 资源类型筛选（可选）
  - `resource_id`: 资源ID筛选（可选）
- **权限**: 需要 `audit:read` 权限

### 获取用户审计日志
- **URL**: `GET /api/v1/audit/users/{user_id}/logs`
- **参数**:
  - `page`: 页码（默认1）
  - `size`: 每页大小（默认20，最大100）
- **权限**: 需要 `audit:read` 权限

### 获取资源审计日志
- **URL**: `GET /api/v1/audit/resources/{resource}/logs`
- **参数**:
  - `resource_id`: 资源ID筛选（可选）
  - `page`: 页码（默认1）
  - `size`: 每页大小（默认20，最大100）
- **权限**: 需要 `audit:read` 权限

## 数据模型

```go
type AuditLog struct {
    ID         string    `json:"id"`          // 唯一标识
    TenantID   string    `json:"tenant_id"`   // 租户ID
    UserID     string    `json:"user_id"`     // 用户ID
    Action     string    `json:"action"`      // 操作类型
    Resource   string    `json:"resource"`    // 资源类型
    ResourceID string    `json:"resource_id"` // 资源ID
    Details    string    `json:"details"`     // 操作详情（JSON）
    Result     string    `json:"result"`      // 结果
    ErrorMsg   string    `json:"error_msg"`   // 错误信息
    CreatedAt  time.Time `json:"created_at"`  // 创建时间
}
```

## 配置说明

### 1. 敏感数据检测
系统会自动检测以下敏感字段：
- password
- secret
- token
- key
- credential

包含这些字段的请求体不会被记录到审计日志中。

### 2. 异步写入
审计日志采用异步写入方式，避免影响主请求的性能。如果审计日志写入失败，会记录到系统日志中，但不会影响主请求的执行。

### 3. 资源ID提取
系统会自动尝试从URL路径中提取UUID格式的资源ID（36个字符，包含4个连字符）。

## 注意事项

1. **性能影响**: 审计中间件采用异步写入，对主请求性能影响很小
2. **存储空间**: 审计日志会持续增长，建议定期清理或归档
3. **权限控制**: 审计日志查询需要相应权限，确保数据安全
4. **多租户**: 所有审计日志都会按租户隔离，确保数据安全
5. **不可修改**: 审计日志一旦创建就不能修改或删除，确保审计的完整性

## 故障排查

### 1. 审计日志没有记录
- 检查是否为写操作（POST/PUT/PATCH/DELETE）
- 检查认证中间件是否正确设置了用户信息
- 检查数据库连接是否正常
- 查看系统日志是否有错误信息

### 2. 查询审计日志失败
- 检查用户是否有 `audit:read` 权限
- 检查租户ID是否正确
- 检查数据库连接是否正常

### 3. 敏感数据被记录
- 检查敏感数据检测逻辑
- 确认字段名是否包含敏感关键词
- 考虑添加新的敏感关键词