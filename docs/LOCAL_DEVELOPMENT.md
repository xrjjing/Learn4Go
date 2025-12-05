# 本地开发指南

本文档详细说明如何在本地环境（不使用 Docker）运行 Learn4Go 项目。

## 🎯 三种本地运行方式

### 方式一：一键启动脚本（推荐）⭐

最简单的方式，自动启动所有服务。

```bash
# 进入项目目录
cd /Users/xrj/GoProject/Learn4Go-1

# 使用内存存储（默认）
./start-local.sh

# 或使用 SQLite
./start-local.sh sqlite

# 或使用 MySQL（需要先启动 Docker MySQL）
docker-compose -f deployments/docker-compose.yml up -d mysql
./start-local.sh mysql
```

**停止服务**：

```bash
./stop-local.sh
# 或按 Ctrl+C
```

**访问地址**：
- 学习门户: http://localhost:8000/portal.html
- 项目首页: http://localhost:8000/index.html
- 项目实战: http://localhost:8000/projects.html
- TODO API: http://localhost:8080/todos
- Gateway: http://localhost:8888/api/todos/todos

**查看日志**：

```bash
# 查看 TODO API 日志
tail -f logs/todoapi.log

# 查看 Gateway 日志
tail -f logs/gateway.log

# 查看前端日志
tail -f logs/frontend.log

# 同时查看所有日志
tail -f logs/*.log
```

---

### 方式二：手动启动（适合调试）

分别在不同终端启动各个服务。

#### 终端 1: 启动 TODO API

```bash
cd /Users/xrj/GoProject/Learn4Go-1

# 选项 A: 内存存储（最简单）
go run ./cmd/todoapi

# 选项 B: SQLite 存储
TODO_STORAGE=sqlite go run ./cmd/todoapi

# 选项 C: MySQL 存储（需要先启动 MySQL）
TODO_STORAGE=mysql \
TODO_DB_HOST=localhost \
TODO_DB_PORT=3306 \
TODO_DB_USER=root \
TODO_DB_PASS=root \
TODO_DB_NAME=learn4go \
go run ./cmd/todoapi
```

#### 终端 2: 启动 API 网关

```bash
cd /Users/xrj/GoProject/Learn4Go-1

GATEWAY_ADDR=:8888 \
TODO_API_URL=http://localhost:8080 \
go run ./examples/gateway/gin
```

#### 终端 3: 启动前端服务器

```bash
cd /Users/xrj/GoProject/Learn4Go-1/web

# 使用 Python（推荐）
python3 -m http.server 8000

# 或使用 Node.js
npx http-server -p 8000

# 或使用 PHP
php -S localhost:8000
```

---

### 方式三：本地代码 + Docker 基础设施

代码在本地运行，但使用 Docker 提供的数据库等基础设施。

#### 步骤 1: 启动基础设施

```bash
cd /Users/xrj/GoProject/Learn4Go-1/deployments

# 只启动 MySQL、Redis、MinIO
docker-compose up -d mysql redis minio

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f mysql
```

#### 步骤 2: 启动应用服务

```bash
cd /Users/xrj/GoProject/Learn4Go-1

# 终端 1: TODO API
TODO_STORAGE=mysql \
TODO_DB_HOST=localhost \
TODO_DB_PORT=3306 \
TODO_DB_USER=root \
TODO_DB_PASS=root \
TODO_DB_NAME=learn4go \
go run ./cmd/todoapi

# 终端 2: Gateway
GATEWAY_ADDR=:8888 \
TODO_API_URL=http://localhost:8080 \
go run ./examples/gateway/gin

# 终端 3: Frontend
cd web && python3 -m http.server 8000
```

#### 停止基础设施

```bash
cd deployments
docker-compose down
```

---

## 📋 环境变量说明

### TODO API 环境变量

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `TODO_STORAGE` | 存储类型 | `memory` | `memory`, `sqlite`, `mysql` |
| `TODO_DB_HOST` | 数据库主机 | `localhost` | `localhost`, `127.0.0.1` |
| `TODO_DB_PORT` | 数据库端口 | `3306` | `3306` |
| `TODO_DB_USER` | 数据库用户 | `root` | `root`, `gouser` |
| `TODO_DB_PASS` | 数据库密码 | - | `root`, `password123` |
| `TODO_DB_NAME` | 数据库名称 | `learn4go` | `learn4go` |

### Gateway 环境变量

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `GATEWAY_ADDR` | 监听地址 | `:8888` | `:8888`, `0.0.0.0:8888` |
| `TODO_API_URL` | TODO API 地址 | `http://localhost:8080` | `http://localhost:8080` |

---

## 🔍 验证服务状态

### 检查服务是否启动

```bash
# TODO API 健康检查
curl http://localhost:8080/healthz

# Gateway 健康检查
curl http://localhost:8888/health

# 前端服务检查
curl http://localhost:8000
```

### 测试 TODO API

```bash
# 获取所有 TODO
curl http://localhost:8080/todos

# 创建 TODO
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"本地测试"}'

# 通过 Gateway 访问
curl http://localhost:8888/api/todos/todos
```

---

## 🐛 常见问题

### 问题 1: 端口被占用

**错误信息**：
```
listen tcp :8080: bind: address already in use
```

**解决方案**：

```bash
# 查找占用端口的进程
lsof -i :8080

# 停止进程
kill <PID>

# 或使用 stop-local.sh
./stop-local.sh
```

### 问题 2: MySQL 连接失败

**错误信息**：
```
dial tcp 127.0.0.1:3306: connect: connection refused
```

**解决方案**：

```bash
# 检查 MySQL 是否启动
docker-compose -f deployments/docker-compose.yml ps mysql

# 启动 MySQL
docker-compose -f deployments/docker-compose.yml up -d mysql

# 查看 MySQL 日志
docker-compose -f deployments/docker-compose.yml logs mysql
```

### 问题 3: 前端无法连接后端

**现象**：前端页面显示服务离线

**解决方案**：

1. 检查后端服务是否启动：
   ```bash
   curl http://localhost:8080/healthz
   curl http://localhost:8888/health
   ```

2. 检查浏览器控制台是否有 CORS 错误

3. 确认前端配置文件 `web/config.js` 中的地址正确

### 问题 4: SQLite 数据库文件权限问题

**错误信息**：
```
unable to open database file
```

**解决方案**：

```bash
# 检查当前目录权限
ls -la

# 删除旧的数据库文件
rm -f todo.db

# 重新启动
TODO_STORAGE=sqlite go run ./cmd/todoapi
```

### 问题 5: Go 依赖下载失败

**错误信息**：
```
go: module ... not found
```

**解决方案**：

```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download

# 整理依赖
go mod tidy
```

---

## 🔧 开发技巧

### 1. 热重载

使用 `air` 实现代码修改后自动重启：

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 在项目根目录运行
air -c .air.toml
```

创建 `.air.toml` 配置文件：

```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/todoapi ./cmd/todoapi"
  bin = "tmp/todoapi"
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor"]
```

### 2. 查看实时日志

```bash
# TODO API 日志
go run ./cmd/todoapi 2>&1 | tee todoapi.log

# Gateway 日志
go run ./examples/gateway/gin 2>&1 | tee gateway.log
```

### 3. 调试模式

```bash
# 使用 delve 调试器
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试 TODO API
dlv debug ./cmd/todoapi

# 在代码中设置断点
(dlv) break main.main
(dlv) continue
```

### 4. 性能分析

```bash
# 启用 pprof
go run ./cmd/todoapi -cpuprofile=cpu.prof

# 分析 CPU 性能
go tool pprof cpu.prof

# 在浏览器中查看
go tool pprof -http=:6060 cpu.prof
```

---

## 📊 端口使用总览

| 服务 | 端口 | 说明 |
|------|------|------|
| TODO API | 8080 | REST API 服务 |
| Gateway | 8888 | API 网关 |
| Frontend | 8000 | 前端静态服务器 |
| MySQL | 3306 | 数据库（Docker） |
| Redis | 6379 | 缓存（Docker） |
| MinIO | 9000 | 对象存储（Docker） |
| MinIO Console | 9001 | MinIO 管理界面（Docker） |

---

## 🎓 学习建议

### 初学者

1. 先使用**内存模式**启动，最简单：
   ```bash
   ./start-local.sh
   ```

2. 访问学习门户，按章节学习

3. 熟悉后再尝试 SQLite 或 MySQL

### 开发者

1. 使用**手动启动**方式，便于调试

2. 在不同终端查看各服务日志

3. 使用 Postman 测试 API

4. 修改代码后重启服务查看效果

### 架构师

1. 使用**本地代码 + Docker 基础设施**方式

2. 研究服务间通信模式

3. 尝试修改 Gateway 路由规则

4. 集成 Redis 缓存和 MinIO 存储

---

## 📚 相关文档

- [README.md](../README.md) - 项目总览
- [API 文档](API.md) - TODO API 接口
- [前端使用指南](FRONTEND.md) - 前端页面说明
- [部署指南](../deployments/README.md) - Docker 部署

---

## 💡 快速命令参考

```bash
# 一键启动（内存模式）
./start-local.sh

# 一键启动（SQLite）
./start-local.sh sqlite

# 一键启动（MySQL）
./start-local.sh mysql

# 停止所有服务
./stop-local.sh

# 查看端口占用
lsof -i :8080
lsof -i :8888
lsof -i :8000

# 测试 API
curl http://localhost:8080/healthz
curl http://localhost:8080/todos

# 查看 Go 进程
ps aux | grep "go run"

# 清理构建缓存
go clean -cache
```

---

**祝你学习愉快！🎉**
