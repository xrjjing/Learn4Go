# API 网关示例

本目录包含两个 API 网关实现，展示不同技术栈的对比：

| 特性 | 标准库版本 | Gin 版本 |
|------|-----------|----------|
| 依赖 | 无 | github.com/gin-gonic/gin |
| 路由 | 手动匹配 | 内置路由器 |
| 中间件 | 自定义实现 | 框架支持 |
| 性能 | 良好 | 优秀（httprouter） |
| 适用场景 | 学习、轻量需求 | 生产环境 |

## 目录结构

```
examples/gateway/
├── stdlib/
│   └── main.go     # 标准库实现
├── gin/
│   └── main.go     # Gin 框架实现
└── README.md
```

## 快速开始

### 1. 启动后端服务

```bash
go run ./cmd/todoapi
```

### 2. 启动网关

**标准库版本：**
```bash
go run ./examples/gateway/stdlib
```

**Gin 版本：**
```bash
# 首次运行需安装依赖
go get github.com/gin-gonic/gin

go run ./examples/gateway/gin
```

### 3. 测试

```bash
# 健康检查
curl http://localhost:8888/health

# 通过网关访问 TODO API（对应后端 /v1/todos）
curl http://localhost:8888/api/v1/todos

# 创建 TODO
curl -X POST http://localhost:8888/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"学习Go网关"}'
```

## 网关功能说明

### 中间件链

```
请求 -> CORS -> 日志 -> 认证 -> [限流] -> 反向代理 -> 后端服务
```

### 已实现的中间件

1. **CORS**: 跨域资源共享
2. **Logger**: 请求日志记录
3. **Auth**: 认证检查（开发模式为警告）
4. **RateLimit**: 限流（仅 Gin 版本）

### 路由映射

| 网关路径 | 后端服务 |
|---------|---------|
| /api/v1/todos/* | localhost:8080 (TODO API 的 /v1/todos) |
| /api/v1/users/* | localhost:8081 |
| /health | 网关自身 |

## 生产环境建议

1. **认证**: 使用 JWT 或 OAuth2
2. **限流**: 使用 Redis 实现分布式限流
3. **熔断**: 集成 hystrix-go 或 sentinel-go
4. **监控**: 添加 Prometheus metrics
5. **追踪**: 集成 OpenTelemetry
