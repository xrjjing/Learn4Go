# TODO API 接口文档

本文档详细说明 TODO REST API 的所有接口。

## 📋 目录

- [基本信息](#基本信息)
- [接口列表](#接口列表)
- [数据模型](#数据模型)
- [错误处理](#错误处理)
- [使用示例](#使用示例)

## 基本信息

### 服务地址

- **Docker 部署**: http://localhost/api/todos
- **本地开发**: http://localhost:8080

### 数据格式

- **请求**: `Content-Type: application/json`
- **响应**: `Content-Type: application/json`

### 存储模式

支持三种存储模式（通过环境变量 `TODO_STORAGE` 配置）：

1. **memory**: 内存存储（默认，重启后数据丢失）
2. **sqlite**: SQLite 文件存储
3. **mysql**: MySQL 数据库存储

## 接口列表

### 认证说明

**重要**: 从 v1.1.0 开始，所有 `/todos*` 接口需要 JWT 认证。

**认证方式**:

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

**获取 Token**: 通过 `/login` 接口登录获取 JWT token。

**Mock 用户**（用于测试）:

- `admin@example.com` / `admin123`
- `user@example.com` / `user123`
- `demo@example.com` / `demo123`

详细认证文档请参考 [JWT 认证系统](JWT认证系统.md)

---

### 1. 用户注册

注册新用户账户。

**请求**

```http
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**参数说明**

| 字段     | 类型   | 必填 | 说明                    |
| -------- | ------ | ---- | ----------------------- |
| email    | string | 是   | 用户邮箱，必须唯一      |
| password | string | 是   | 用户密码，建议 8 位以上 |

**响应**

```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": 4,
  "email": "user@example.com",
  "created_at": "2025-12-05T14:00:00Z"
}
```

**错误响应**

```json
HTTP/1.1 409 Conflict
{
  "error": "email already exists"
}
```

**示例**

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","password":"pass123"}'
```

---

### 2. 用户登录

用户登录获取 JWT access 与 refresh token，并触发登录失败次数限制。

**请求**

```http
POST /login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "admin123"
}
```

**参数说明**

| 字段     | 类型   | 必填 | 说明     |
| -------- | ------ | ---- | -------- |
| email    | string | 是   | 用户邮箱 |
| password | string | 是   | 用户密码 |

**响应**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400,
  "refresh_token": "QWxhZGRpbjpPcGVuU2VzYW1l",
  "refresh_expires_in": 604800,
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

**错误响应**

```json
HTTP/1.1 401 Unauthorized
{
  "error": "invalid credentials"
}

HTTP/1.1 429 Too Many Requests
Retry-After: 600
{
  "error": "account temporarily locked"
}
```

**示例**

```bash
# 登录并保存token
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.token')

echo "Token: $TOKEN"
```

---

### 3. 刷新令牌

使用 refresh token 旋转获取新的 access/refresh。

```http
POST /refresh
Content-Type: application/json

{
  "refresh_token": "<refresh-from-login>"
}
```

**响应**

```json
{
  "token": "new-access",
  "expires_in": 86400,
  "refresh_token": "new-refresh",
  "refresh_expires_in": 604800
}
```

---

### 3. 健康检查

检查服务是否正常运行。

**请求**

```http
GET /healthz
```

**响应**

```http
HTTP/1.1 200 OK
Content-Type: text/plain

ok
```

**示例**

```bash
curl http://localhost:8080/healthz
```

---

### 4. 获取所有 TODO

获取所有待办事项列表。**需要认证**。

**请求**

```http
GET /todos
Authorization: Bearer YOUR_JWT_TOKEN
```

**响应**

```json
HTTP/1.1 200 OK
Content-Type: application/json

[
  {
    "id": 1,
    "title": "学习 Go 语言基础",
    "done": false,
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "title": "完成 TODO API 项目",
    "done": true,
    "created_at": "2024-01-15T11:00:00Z"
  }
]
```

**错误响应**

```json
HTTP/1.1 401 Unauthorized
{
  "error": "authorization required"
}
```

**示例**

```bash
# 使用token访问
curl http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN"

# 通过网关访问
curl http://localhost:8888/api/todos/todos \
  -H "Authorization: Bearer $TOKEN"
```

---

### 5. 创建 TODO

创建一个新的待办事项。**需要认证**。

**请求**

```http
POST /todos
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
  "title": "待办事项标题"
}
```

**参数说明**

| 字段  | 类型   | 必填 | 说明                   |
| ----- | ------ | ---- | ---------------------- |
| title | string | 是   | 待办事项标题，不能为空 |

**响应**

```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": 3,
  "title": "学习 Gin 框架",
  "done": false,
  "created_at": "2024-01-15T12:00:00Z"
}
```

**错误响应**

```json
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": "title required"
}
```

**示例**

```bash
# 创建 TODO（需要token）
curl -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 Go 并发编程"}'

# 通过网关创建
curl -X POST http://localhost:8888/api/todos/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"部署到 Docker"}'
```

---

### 6. 更新 TODO 状态

更新待办事项的完成状态。**需要认证**。

**请求**

```http
PUT /todos/{id}
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
  "done": true
}
```

**路径参数**

| 参数 | 类型    | 说明        |
| ---- | ------- | ----------- |
| id   | integer | 待办事项 ID |

**请求体参数**

| 字段 | 类型    | 必填 | 说明                      |
| ---- | ------- | ---- | ------------------------- |
| done | boolean | 是   | 完成状态，true 表示已完成 |

**响应**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": 1,
  "title": "学习 Go 语言基础",
  "done": true,
  "created_at": "2024-01-15T10:30:00Z"
}
```

**错误响应**

```json
HTTP/1.1 404 Not Found
Content-Type: application/json

{
  "error": "not found"
}
```

**示例**

```bash
# 标记为已完成（需要token）
curl -X PUT http://localhost:8080/todos/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"done":true}'

# 标记为未完成
curl -X PUT http://localhost:8080/todos/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"done":false}'
```

---

### 7. 删除 TODO

删除指定的待办事项。**需要认证**。

**请求**

```http
DELETE /todos/{id}
Authorization: Bearer YOUR_JWT_TOKEN
```

**路径参数**

| 参数 | 类型    | 说明        |
| ---- | ------- | ----------- |
| id   | integer | 待办事项 ID |

**响应**

```http
HTTP/1.1 204 No Content
```

**错误响应**

```json
HTTP/1.1 404 Not Found
Content-Type: application/json

{
  "error": "not found"
}
```

**示例**

```bash
# 删除 TODO（需要token）
curl -X DELETE http://localhost:8080/todos/1 \
  -H "Authorization: Bearer $TOKEN"

# 通过网关删除
curl -X DELETE http://localhost:8888/api/todos/todos/1 \
  -H "Authorization: Bearer $TOKEN"
```

## 数据模型

### Todo 对象

| 字段       | 类型    | 说明                    |
| ---------- | ------- | ----------------------- |
| id         | integer | 唯一标识符，自动生成    |
| title      | string  | 待办事项标题            |
| done       | boolean | 完成状态，默认 false    |
| created_at | string  | 创建时间，ISO 8601 格式 |

**示例**

```json
{
  "id": 1,
  "title": "学习 Go 语言",
  "done": false,
  "created_at": "2024-01-15T10:30:00Z"
}
```

## 错误处理

### 错误响应格式

所有错误响应都遵循统一格式：

```json
{
  "error": "错误描述信息"
}
```

### HTTP 状态码

| 状态码                    | 说明                 | 场景                        |
| ------------------------- | -------------------- | --------------------------- |
| 200 OK                    | 请求成功             | GET、PUT 成功               |
| 201 Created               | 资源创建成功         | POST 成功                   |
| 204 No Content            | 请求成功，无返回内容 | DELETE 成功                 |
| 400 Bad Request           | 请求参数错误         | 缺少必填字段、JSON 格式错误 |
| 404 Not Found             | 资源不存在           | ID 不存在                   |
| 405 Method Not Allowed    | 方法不允许           | 使用了不支持的 HTTP 方法    |
| 500 Internal Server Error | 服务器内部错误       | 数据库错误等                |

### 常见错误

#### 1. 缺少标题

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":""}'
```

响应：

```json
{
  "error": "title required"
}
```

#### 2. JSON 格式错误

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{invalid json}'
```

响应：

```json
{
  "error": "invalid json"
}
```

#### 3. TODO 不存在

```bash
curl -X PUT http://localhost:8080/todos/999 \
  -H "Content-Type: application/json" \
  -d '{"done":true}'
```

响应：

```json
{
  "error": "not found"
}
```

#### 4. 数据库错误

当数据库连接失败或查询出错时：

```json
{
  "error": "internal error"
}
```

## 使用示例

### 完整工作流（含认证）

```bash
# 1. 检查服务健康状态
curl http://localhost:8080/healthz

# 2. 用户登录获取token
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.token')

echo "Token: $TOKEN"

# 3. 获取所有 TODO（初始为空）
curl http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN"

# 4. 创建第一个 TODO
curl -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 Go 语言基础"}'

# 5. 创建第二个 TODO
curl -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"完成 TODO API 项目"}'

# 6. 查看所有 TODO
curl http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN"

# 7. 标记第一个 TODO 为已完成
curl -X PUT http://localhost:8080/todos/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"done":true}'

# 8. 删除第二个 TODO
curl -X DELETE http://localhost:8080/todos/2 \
  -H "Authorization: Bearer $TOKEN"

# 9. 再次查看所有 TODO
curl http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN"
```

### 使用 jq 格式化输出

```bash
# 安装 jq
# macOS: brew install jq
# Ubuntu: sudo apt-get install jq

# 格式化输出
curl -s http://localhost:8080/todos | jq .

# 只显示标题
curl -s http://localhost:8080/todos | jq '.[].title'

# 只显示未完成的 TODO
curl -s http://localhost:8080/todos | jq '.[] | select(.done == false)'
```

### 使用 Postman

1. 导入以下 Collection：

```json
{
  "info": {
    "name": "Learn4Go TODO API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/healthz"
      }
    },
    {
      "name": "Get All TODOs",
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/todos"
      }
    },
    {
      "name": "Create TODO",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\"title\":\"学习 Go 语言\"}"
        },
        "url": "http://localhost:8080/todos"
      }
    },
    {
      "name": "Update TODO",
      "request": {
        "method": "PUT",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\"done\":true}"
        },
        "url": "http://localhost:8080/todos/1"
      }
    },
    {
      "name": "Delete TODO",
      "request": {
        "method": "DELETE",
        "url": "http://localhost:8080/todos/1"
      }
    }
  ]
}
```

## 环境变量配置

### TODO API 配置

| 变量名       | 说明       | 默认值    | 示例                  |
| ------------ | ---------- | --------- | --------------------- |
| TODO_STORAGE | 存储类型   | memory    | memory, sqlite, mysql |
| TODO_DB_HOST | 数据库主机 | localhost | mysql, 127.0.0.1      |
| TODO_DB_PORT | 数据库端口 | 3306      | 3306                  |
| TODO_DB_USER | 数据库用户 | root      | root, gouser          |
| TODO_DB_PASS | 数据库密码 | -         | password123           |
| TODO_DB_NAME | 数据库名称 | learn4go  | learn4go              |

### 使用示例

```bash
# 使用内存存储
go run ./cmd/todoapi

# 使用 SQLite
TODO_STORAGE=sqlite go run ./cmd/todoapi

# 使用 MySQL
TODO_STORAGE=mysql \
TODO_DB_HOST=localhost \
TODO_DB_PORT=3306 \
TODO_DB_USER=root \
TODO_DB_PASS=root \
TODO_DB_NAME=learn4go \
go run ./cmd/todoapi
```

## 性能考虑

### 并发安全

- **内存存储**: 使用 `sync.Mutex` 保证并发安全
- **数据库存储**: 依赖数据库事务保证一致性

### 连接池

MySQL 模式下，GORM 自动管理连接池：

- 默认最大空闲连接：2
- 默认最大打开连接：无限制
- 连接最大生命周期：无限制

### 性能优化建议

1. **使用数据库索引**: 在 `id` 字段上已有主键索引
2. **启用查询缓存**: 可集成 Redis 缓存热点数据
3. **批量操作**: 对于大量数据，考虑批量插入/更新
4. **分页查询**: 当数据量大时，添加分页参数

## 安全考虑

### 当前实现

- ✅ 输入验证（标题非空）
- ✅ JSON 格式验证
- ✅ SQL 注入防护（GORM 参数化查询）
- ✅ 错误信息脱敏（不暴露内部错误）
- ✅ JWT 认证（v1.1.0+）
- ✅ bcrypt 密码加密
- ✅ 基于内存的速率限制

### 生产环境建议

- [ ] 添加授权（RBAC）
- [ ] 添加 HTTPS/TLS
- [ ] 添加 CORS 配置
- [ ] 添加请求日志
- [ ] 添加审计日志
- [ ] Token 刷新机制
- [ ] 登录失败次数限制
- [ ] 密码强度验证

## 相关文档

- [README.md](../README.md) - 项目总览
- [认证系统文档](JWT认证系统.md) - JWT 认证详细说明
- [前端使用指南](前端使用指南.md) - 前端页面说明
- [部署指南](../deployments/README.md) - Docker 部署
- [项目计划](../plan.md) - 后续开发计划

## 问题反馈

如有问题或建议，欢迎通过 Issue 反馈。
