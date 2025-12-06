# scripts 目录脚本索引

本目录提供了一组 **极简包装脚本**，统一入口，方便快速执行常用的启动 / 停止 / 测试 / 部署操作。

> 详细的链路说明与测试步骤请参考：`docs/前后端链路与测试指南.md`

---

## 一、开发阶段脚本（直连 TODO API）

### 1. 启动本地开发环境

- 脚本：`scripts/dev-start.sh`
- 等价命令：`./start-local.sh`
- 功能：
  - 启动 TODO API（默认内存模式）
  - 启动 Gin 网关示例（端口 8888）
  - 启动前端静态服务器（端口 8000）
  - 做一次健康检查，并打印访问入口：
    - 学习门户：`http://localhost:8000/portal.html`
    - TODO API：`http://localhost:8080/v1/todos`
    - Gateway：`http://localhost:8888/api/v1/todos`

使用示例：

```bash
./scripts/dev-start.sh            # 内存模式
./scripts/dev-start.sh sqlite     # SQLite 模式
./scripts/dev-start.sh mysql      # MySQL 模式（需先启动数据库）
```

### 2. 停止本地开发环境

- 脚本：`scripts/dev-stop.sh`
- 等价命令：`./stop-local.sh`
- 功能：
  - 读取 `logs/.pids` 中记录的 PID
  - 停止 TODO API / Gateway / Frontend 等本地进程

使用示例：

```bash
./scripts/dev-stop.sh
```

---

## 二、测试脚本

### 1. TODO API 冒烟测试

- 脚本：`scripts/test-smoke.sh`
- 底层调用：项目根目录的 `./test-smoke.sh`
- 功能：
  1. 检查 `http://127.0.0.1:8080/healthz`
  2. 使用 `admin@example.com / admin123` 登录，获取 token
  3. 带 token 访问 `GET /v1/todos`
  4. 如检测到 Gateway 监听在 `:8888`，再访问 `GET /api/v1/todos`

使用示例：

```bash
./scripts/test-smoke.sh
```

> 提示：请先在另一个终端执行 `go run ./cmd/todoapi` 或 `./scripts/dev-start.sh`。

### 2. 管理后台 API 集成测试

- 脚本：`scripts/test-admin.sh`
- 底层调用：项目根目录的 `./test-admin-api.sh`
- 功能：
  - 检查 `/healthz`
  - 登录 admin 并测试：
    - `GET /v1/me`
    - `GET /v1/users`
    - `GET /v1/rbac/roles`
    - `GET /v1/rbac/permissions`
    - `POST /v1/users`
    - `POST /v1/logout`

使用示例：

```bash
./scripts/test-admin.sh
```

> 建议配合 `web/login.html` / `web/admin.html` 一起使用，做端到端验证。

---

## 三、Docker 部署脚本

### 1. 启动 Docker 环境

- 脚本：`scripts/deploy-up.sh`
- 等价命令：在 `deployments` 目录执行 `docker-compose up -d`
- 功能：
  - 启动 todoapi、gateway、nginx 以及数据库、Redis、MinIO 等服务
  - 通过 Nginx 对外暴露：
    - 前端门户：`http://localhost`
    - TODO API：`http://localhost/api/v1/todos`

使用示例：

```bash
./scripts/deploy-up.sh
```

### 2. 停止 Docker 环境（保留数据）

- 脚本：`scripts/deploy-down.sh`
- 等价命令：在 `deployments` 目录执行 `docker-compose down`
- 功能：
  - 停止所有容器，但保留数据卷（数据库数据不会丢失）

使用示例：

```bash
./scripts/deploy-down.sh
```

> 如需同时删除数据卷，可在 `deployments` 目录手动执行：  
> `docker-compose down -v`（注意：这是高风险操作，会清除数据）。

---

## 四、快速记忆表

| 脚本                     | 作用                             | 典型场景           |
|--------------------------|----------------------------------|--------------------|
| `scripts/dev-start.sh`   | 启动本地开发环境                 | 日常开发 / 联调    |
| `scripts/dev-stop.sh`    | 停止本地开发环境                 | 结束开发           |
| `scripts/test-smoke.sh`  | TODO API 冒烟测试（登录 + 列表） | 核心链路自检       |
| `scripts/test-admin.sh`  | 管理后台 API 集成测试            | 管理后台功能验证   |
| `scripts/deploy-up.sh`   | Docker 环境启动                  | 部署 / 架构演示    |
| `scripts/deploy-down.sh` | Docker 环境停止（保留数据）      | 停止部署环境       |

> 更详细的链路说明和 curl 示例，请配合 `docs/前后端链路与测试指南.md` 一起阅读使用。

