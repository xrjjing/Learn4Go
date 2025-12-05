# JWT 认证系统文档

## 概述

本项目实现了基于 JWT (JSON Web Token) 的用户认证系统，提供了完整的用户注册、登录和API访问控制功能。

## 技术栈

- **JWT**: `github.com/golang-jwt/jwt/v5` - JWT token生成和验证
- **密码加密**: `golang.org/x/crypto/bcrypt` - bcrypt密码哈希
- **签名算法**: HS256 (HMAC-SHA256)
- **Token有效期**: 24小时（可配置）

## 核心功能

### 1. 用户注册 (`POST /register`)

**请求示例**：
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "password123"
  }'
```

**响应示例**：
```json
{
  "id": 4,
  "email": "newuser@example.com",
  "created_at": "2025-12-05T14:00:00Z"
}
```

**特性**：
- 邮箱唯一性验证
- bcrypt密码加密（cost=10）
- 自动生成用户ID

### 2. 用户登录 (`POST /login`)

**请求示例**：
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

**响应示例**：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400,
  "user": {
    "id": 1,
    "email": "admin@example.com"
  }
}
```

**特性**：
- 密码验证使用bcrypt
- 返回JWT token和过期时间
- 登录失败返回统一错误信息（防止用户枚举）

### 3. 受保护的API访问

所有 `/todos*` 路径需要JWT认证。

**请求示例**：
```bash
# 获取TODO列表
curl http://localhost:8080/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 创建TODO
curl -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "新任务"}'
```

**未认证访问返回**：
```json
{
  "error": "authorization required"
}
```

**Token无效或过期返回**：
```json
{
  "error": "invalid or expired token"
}
```

## Mock用户数据

系统预置了3个测试用户（内存存储）：

| 邮箱 | 密码 | 用途 |
|------|------|------|
| admin@example.com | admin123 | 管理员账户 |
| user@example.com | user123 | 普通用户 |
| demo@example.com | demo123 | 演示账户 |

## 架构设计

### 1. 认证中间件

```go
// authMiddleware 保护 /todos* 路径
func (s *Server) authMiddleware(next http.Handler) http.Handler {
    // 公开路径：/, /healthz, /register, /login
    // 受保护路径：/todos*
}
```

**中间件链**：
```
请求 → 日志中间件 → 限流中间件 → 认证中间件 → 业务处理器
```

### 2. JWT Manager

```go
type JWTManager struct {
    cfg JWTConfig
}

// 生成token
func (m *JWTManager) Generate(userID uint) (string, error)

// 验证token
func (m *JWTManager) Parse(tokenStr string) (uint, error)
```

**Token Claims**：
- `sub`: 用户ID
- `iat`: 签发时间
- `exp`: 过期时间

### 3. 用户存储接口

```go
type UserStore interface {
    Create(ctx context.Context, email, passwordHash string) (User, error)
    FindByEmail(ctx context.Context, email string) (User, error)
    FindByID(ctx context.Context, id uint) (User, error)
}
```

**实现**：
- `MemoryUserStore`: 内存存储（开发/测试）
- `DBStore`: 数据库存储（生产环境，待实现）

## 安全特性

### 1. 密码安全
- ✅ bcrypt加密（cost=10）
- ✅ 密码不在响应中返回（`json:"-"`）
- ✅ 登录失败统一错误信息

### 2. Token安全
- ✅ HS256签名算法
- ✅ 24小时过期时间
- ✅ 签名密钥可配置
- ⚠️ 生产环境需更换默认密钥

### 3. API保护
- ✅ 所有TODO API需要认证
- ✅ Token验证失败返回401
- ✅ 用户ID存储在请求上下文中

## 配置选项

```go
// 创建服务器时配置JWT
s := todo.NewServer(store,
    todo.WithJWT("your-secret-key", 24*time.Hour),
    todo.WithUserStore(userStore),
    todo.WithRateLimiter(rateLimiter),
)
```

**环境变量**（建议）：
```bash
export JWT_SECRET="your-production-secret-key"
export JWT_TTL="24h"
```

## 使用示例

### 完整认证流程

```bash
# 1. 注册新用户
REGISTER_RESP=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}')
echo "注册响应: $REGISTER_RESP"

# 2. 登录获取token
LOGIN_RESP=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}')
TOKEN=$(echo $LOGIN_RESP | jq -r '.token')
echo "Token: $TOKEN"

# 3. 使用token访问API
curl -s http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN"

# 4. 创建TODO
curl -s -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"我的第一个任务"}'
```

## 错误处理

| HTTP状态码 | 错误信息 | 说明 |
|-----------|---------|------|
| 400 | invalid json | 请求体格式错误 |
| 400 | email and password required | 缺少必填字段 |
| 401 | authorization required | 缺少Authorization头 |
| 401 | invalid or expired token | Token无效或过期 |
| 401 | invalid credentials | 邮箱或密码错误 |
| 409 | email already exists | 邮箱已被注册 |
| 500 | internal error | 服务器内部错误 |

## 后续优化

### 短期（已规划）
- [ ] Token刷新机制（Refresh Token）
- [ ] 登录失败次数限制
- [ ] 密码强度验证
- [ ] 邮箱格式验证

### 中期
- [ ] RBAC权限控制
- [ ] OAuth2集成（Google/GitHub）
- [ ] 双因素认证（2FA）
- [ ] 会话管理（黑名单）

### 长期
- [ ] 分布式Session（Redis）
- [ ] 单点登录（SSO）
- [ ] 审计日志
- [ ] 安全事件告警

## Java开发者对比

| 概念 | Java (Spring Security) | Go (本项目) |
|------|----------------------|------------|
| 认证过滤器 | `OncePerRequestFilter` | `authMiddleware` |
| 密码加密 | `BCryptPasswordEncoder` | `bcrypt.GenerateFromPassword` |
| JWT库 | `jjwt` | `golang-jwt/jwt` |
| 用户存储 | `UserDetailsService` | `UserStore` interface |
| 上下文传递 | `SecurityContextHolder` | `context.WithValue` |
| 配置 | `@EnableWebSecurity` | `WithJWT` Option |

**关键差异**：
1. Go使用中间件链而非过滤器链
2. Go的接口是隐式实现，更灵活
3. Go使用context传递用户信息，而非ThreadLocal
4. Go的错误处理是显式的（`if err != nil`）

## 测试

```bash
# 运行认证相关测试
go test ./internal/todo -run TestAuth

# 测试JWT生成和验证
go test ./internal/todo -run TestJWT

# 测试密码加密
go test ./internal/todo -run TestPassword
```

## 参考资料

- [JWT官方文档](https://jwt.io/)
- [bcrypt算法](https://en.wikipedia.org/wiki/Bcrypt)
- [OWASP认证备忘单](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [Go Context最佳实践](https://go.dev/blog/context)

---

**最后更新**: 2025-12-05
**版本**: v1.0.0
