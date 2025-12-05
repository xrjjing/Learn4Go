# RBAC 权限控制系统文档

## 概述

本项目实现了基于角色的访问控制（Role-Based Access Control, RBAC），提供细粒度的权限管理功能。

## 角色定义

系统预定义了三种角色：

| 角色 | 代码 | 权限范围 | 使用场景 |
|------|------|---------|---------|
| **Admin** | `admin` | 完全权限 | 系统管理员，可以查看和操作所有用户的TODO |
| **User** | `user` | 自己的资源 | 普通用户，只能操作自己创建的TODO |
| **Guest** | `guest` | 只读权限 | 访客用户，可以查看所有TODO但不能修改 |

## 权限矩阵

### TODO资源权限

| 角色 | 创建 | 读取 | 更新 | 删除 | 约束 |
|------|------|------|------|------|------|
| Admin | ✅ | ✅ | ✅ | ✅ | 无限制 |
| User | ✅ | ✅ | ✅ | ✅ | 仅自己的TODO |
| Guest | ❌ | ✅ | ❌ | ❌ | 只读所有TODO |

## 架构设计

### 1. 数据模型

```go
// User 用户实体
type User struct {
    ID           uint
    Email        string
    PasswordHash string
    Role         Role      // admin/user/guest
    CreatedAt    time.Time
}

// Todo 待办实体
type Todo struct {
    ID        int
    UserID    uint      // 所属用户ID
    Title     string
    Done      bool
    CreatedAt time.Time
}
```

### 2. RBAC管理器

```go
type RBACManager struct {
    permissions map[Role][]Permission
}

// 检查权限
func (m *RBACManager) CheckPermission(role Role, resource Resource, action Action) bool
```

### 3. 中间件链

```
请求 → 日志 → 限流 → JWT认证 → RBAC授权 → 业务处理
```

## 使用示例

### 1. Admin用户（完全权限）

```bash
# 登录
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.token')

# 创建TODO
curl -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"管理员任务"}'

# 查看所有TODO（包括其他用户的）
curl http://localhost:8080/todos \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 删除任何TODO
curl -X DELETE http://localhost:8080/todos/123 \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### 2. User用户（自己的资源）

```bash
# 登录
USER_TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"user123"}' \
  | jq -r '.token')

# 创建TODO（自动关联到当前用户）
curl -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"我的任务"}'

# 查看TODO（只能看到自己的）
curl http://localhost:8080/todos \
  -H "Authorization: Bearer $USER_TOKEN"

# 尝试删除其他用户的TODO（会被拒绝）
curl -X DELETE http://localhost:8080/todos/999 \
  -H "Authorization: Bearer $USER_TOKEN"
# 返回: {"error":"you don't own this resource"}
```

### 3. Guest用户（只读）

```bash
# 登录
GUEST_TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"demo@example.com","password":"demo123"}' \
  | jq -r '.token')

# 查看所有TODO（可以看到）
curl http://localhost:8080/todos \
  -H "Authorization: Bearer $GUEST_TOKEN"

# 尝试创建TODO（会被拒绝）
curl -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $GUEST_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"访客任务"}'
# 返回: {"error":"insufficient permissions"}
```

## 权限验证流程

### 1. 基本权限检查

```
1. 从JWT中提取用户ID
2. 查询用户角色
3. 解析请求的资源和操作
4. 检查角色是否有该权限
5. 通过则继续，否则返回403
```

### 2. 所有权验证（User角色）

```
1. 基本权限检查通过
2. 如果是User角色且操作是read/update/delete
3. 从URL提取TODO ID
4. 查询TODO的UserID
5. 比较TODO.UserID与当前用户ID
6. 匹配则通过，否则返回403
```

## 错误响应

| HTTP状态码 | 错误信息 | 说明 |
|-----------|---------|------|
| 401 | authorization required | 未提供JWT token |
| 401 | invalid or expired token | Token无效或过期 |
| 401 | user not found | 用户不存在 |
| 403 | insufficient permissions | 角色权限不足 |
| 403 | you don't own this resource | 不是资源所有者 |

## 数据库表结构

### roles 表

```sql
CREATE TABLE roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### user_roles 表

```sql
CREATE TABLE user_roles (
    user_id INT NOT NULL,
    role_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);
```

### todos 表更新

```sql
ALTER TABLE todos ADD COLUMN user_id INT NOT NULL DEFAULT 1;
ALTER TABLE todos ADD INDEX idx_user_id (user_id);
```

## 配置选项

### 默认角色

新注册用户默认角色为 `user`（普通用户）。

### 角色分配

```go
// 创建用户时指定角色
user, err := userStore.CreateWithRole(ctx, email, hash, RoleUser)

// 更新用户角色（管理员操作）
err := userStore.UpdateRole(ctx, userID, RoleAdmin)
```

## 扩展性

### 1. 添加新角色

```go
const (
    RoleAdmin    Role = "admin"
    RoleUser     Role = "user"
    RoleGuest    Role = "guest"
    RoleModerator Role = "moderator" // 新增角色
)

// 在RBACManager中添加权限
permissions := map[Role][]Permission{
    RoleModerator: {
        {ResourceTodos, ActionRead},
        {ResourceTodos, ActionUpdate},
    },
}
```

### 2. 添加新资源

```go
const (
    ResourceTodos    Resource = "todos"
    ResourceComments Resource = "comments" // 新增资源
)

// 定义新资源的权限
permissions := map[Role][]Permission{
    RoleAdmin: {
        {ResourceComments, ActionCreate},
        {ResourceComments, ActionRead},
        {ResourceComments, ActionUpdate},
        {ResourceComments, ActionDelete},
    },
}
```

### 3. 动态权限（未来）

当前权限规则硬编码在代码中，未来可以：
- 从数据库加载权限配置
- 使用Casbin等权限框架
- 支持动态权限分配

## Java开发者对比

| 概念 | Java (Spring Security) | Go (本项目) |
|------|----------------------|------------|
| 角色定义 | `@RolesAllowed("ADMIN")` | `Role` 常量 |
| 权限检查 | `hasRole()`, `hasAuthority()` | `CheckPermission()` |
| 方法级保护 | `@PreAuthorize` | RBAC中间件 |
| 角色存储 | `GrantedAuthority` | `User.Role` |
| 权限管理 | `AccessDecisionManager` | `RBACManager` |
| 所有权验证 | `@PostAuthorize` | 中间件内验证 |

**关键差异**：
1. Go使用中间件而非注解
2. Go的权限检查是显式的函数调用
3. Go没有AOP，需要手动集成中间件
4. Go的错误处理是显式的

## 测试

### 单元测试

```bash
# 测试RBAC管理器
go test ./internal/todo -run TestRBAC

# 测试权限验证
go test ./internal/todo -run TestPermission
```

### 集成测试

```bash
# 运行完整的RBAC测试脚本
./examples/rbac-demo.sh
```

## 安全考虑

### 当前实现

- ✅ 角色基于权限验证
- ✅ 所有权验证（User角色）
- ✅ 默认拒绝策略
- ✅ 权限不足返回403

### 生产环境建议

- [ ] 审计日志（记录权限变更）
- [ ] 权限缓存（Redis）
- [ ] 动态权限配置
- [ ] 细粒度权限（字段级）
- [ ] 权限继承
- [ ] 临时权限授予

## 故障排查

### 问题1：403 Forbidden

**原因**：用户角色权限不足

**解决**：
1. 检查用户角色：`SELECT role FROM users WHERE id = ?`
2. 确认角色权限配置
3. 查看日志中的权限检查信息

### 问题2：you don't own this resource

**原因**：User尝试操作其他用户的TODO

**解决**：
1. 确认TODO的user_id
2. 确认当前用户ID
3. 使用Admin账户操作

### 问题3：角色未生效

**原因**：JWT中未包含角色信息或角色未正确设置

**解决**：
1. 重新登录获取新token
2. 检查User表的role字段
3. 确认mock数据中的角色设置

## 相关文档

- [JWT认证文档](AUTH.md) - JWT认证系统
- [API文档](API.md) - API接口说明
- [变更日志](CHANGELOG.md) - 版本更新记录

---

**最后更新**: 2025-12-05
**版本**: v1.2.0
