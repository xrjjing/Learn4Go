# Docker 部署指南

本目录包含 Learn4Go 项目的 Docker 容器化部署配置。

## 架构概览

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ :80
       ▼
┌─────────────┐
│   Nginx     │ (前端静态文件 + 反向代理)
└──────┬──────┘
       │ /api/*
       ▼
┌─────────────┐
│   Gateway   │ :8888 (API 网关)
└──────┬──────┘
       │
       ├─────► TODO API :8080 (REST API)
       │
       ├─────► MySQL :3306 (数据库)
       │
       ├─────► Redis :6379 (缓存)
       │
       └─────► MinIO :9000 (对象存储)
```

## 服务列表

| 服务名 | 端口 | 说明 |
|--------|------|------|
| frontend | 80 | Nginx 前端服务器 |
| gateway | 8888 | Gin API 网关 |
| todoapi | 8080 | TODO REST API |
| mysql | 3306 | MySQL 数据库 |
| redis | 6379 | Redis 缓存 |
| minio | 9000, 9001 | MinIO 对象存储 |

## 快速开始

### 1. 构建并启动所有服务

```bash
cd deployments
docker-compose up -d
```

### 2. 查看服务状态

```bash
docker-compose ps
```

### 3. 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f todoapi
docker-compose logs -f gateway
```

### 4. 访问服务

- **前端门户**: http://localhost
- **TODO API**: http://localhost/api/todos
- **API 网关**: http://localhost/api
- **MinIO 控制台**: http://localhost:9001 (用户名: minioadmin, 密码: minioadmin)

### 5. 停止服务

```bash
docker-compose down
```

### 6. 清理数据（包括数据库和存储卷）

```bash
docker-compose down -v
```

## 环境变量配置

### MySQL 配置

```bash
MYSQL_ROOT_PASSWORD=rootpass123
MYSQL_DATABASE=learn4go
MYSQL_USER=gouser
MYSQL_PASSWORD=gopass123
```

### Redis 配置

```bash
REDIS_PASSWORD=redispass123
```

### MinIO 配置

```bash
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
```

### TODO API 配置

```bash
TODO_STORAGE=mysql
TODO_DB_HOST=mysql
TODO_DB_PORT=3306
TODO_DB_USER=gouser
TODO_DB_PASS=gopass123
TODO_DB_NAME=learn4go
```

### Gateway 配置

```bash
GATEWAY_ADDR=:8888
TODO_API_URL=http://todoapi:8080
```

## 开发模式

如果需要在本地开发时使用 Docker 数据库，但运行本地代码：

```bash
# 仅启动基础设施服务
docker-compose up -d mysql redis minio

# 本地运行 TODO API
cd ..
TODO_STORAGE=mysql \
TODO_DB_HOST=localhost \
TODO_DB_PORT=3306 \
TODO_DB_USER=gouser \
TODO_DB_PASS=gopass123 \
TODO_DB_NAME=learn4go \
go run ./cmd/todoapi

# 本地运行 Gateway
GATEWAY_ADDR=:8888 \
TODO_API_URL=http://localhost:8080 \
go run ./examples/gateway/gin
```

## 健康检查

所有服务都配置了健康检查：

- **todoapi**: `GET /healthz`
- **gateway**: `GET /health`
- **mysql**: `mysqladmin ping`
- **redis**: `redis-cli ping`
- **minio**: `curl -f http://localhost:9000/minio/health/live`

## 数据持久化

以下数据会持久化到 Docker 卷：

- `mysql_data`: MySQL 数据库文件
- `redis_data`: Redis 持久化数据
- `minio_data`: MinIO 对象存储数据

## 故障排查

### 服务无法启动

```bash
# 检查容器状态
docker-compose ps

# 查看详细日志
docker-compose logs <service-name>

# 重启特定服务
docker-compose restart <service-name>
```

### 数据库连接失败

```bash
# 进入 MySQL 容器
docker-compose exec mysql mysql -ugouser -pgopass123 learn4go

# 检查数据库是否创建
SHOW DATABASES;
```

### Redis 连接失败

```bash
# 进入 Redis 容器
docker-compose exec redis redis-cli -a redispass123

# 测试连接
PING
```

### MinIO 连接失败

```bash
# 检查 MinIO 日志
docker-compose logs minio

# 访问 MinIO 控制台
open http://localhost:9001
```

## 生产部署建议

1. **修改默认密码**: 在 `docker-compose.yml` 中修改所有默认密码
2. **使用环境变量文件**: 创建 `.env` 文件存储敏感信息
3. **配置 HTTPS**: 在 Nginx 中配置 SSL 证书
4. **资源限制**: 为每个服务配置 CPU 和内存限制
5. **日志管理**: 配置日志轮转和集中日志收集
6. **监控告警**: 集成 Prometheus + Grafana 监控
7. **备份策略**: 定期备份 MySQL 和 MinIO 数据

## 扩展阅读

- [Docker Compose 文档](https://docs.docker.com/compose/)
- [Nginx 反向代理配置](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
- [MySQL Docker 镜像](https://hub.docker.com/_/mysql)
- [Redis Docker 镜像](https://hub.docker.com/_/redis)
- [MinIO Docker 部署](https://min.io/docs/minio/container/index.html)
