# Learn4Go - Go 语言学习与微服务实战项目

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

面向 Java 开发者的 Go 语言系统化学习项目，涵盖**语言基础 → 框架生态 → 微服务架构 → Docker 部署**完整链路。

## ✨ 项目特色

- 📚 **系统化学习路径**: 从基础语法到微服务架构，循序渐进
- 🔄 **Java 对照学习**: 提供 Java vs Go 速查表，快速上手
- 🏗️ **生产级架构**: 完整的微服务技术栈（Gin + gRPC + MySQL + Redis + MinIO）
- 🐳 **一键部署**: Docker Compose 编排，开箱即用
- 🎨 **前端集成**: 学习门户 + TODO 应用演示
- 🔐 **安全示例**: JWT + RBAC + Refresh Token + 登录失败限制示范

## 🚀 快速开始

### 方式一：Docker 部署（推荐）

```bash
# 1. 克隆项目
git clone <your-repo-url>
cd Learn4Go-1

# 2. 启动所有服务
cd deployments
docker-compose up -d

# 3. 访问服务
# 前端门户: http://localhost
# TODO API: http://localhost/api/todos
# MinIO 控制台: http://localhost:9001
```

### 方式二：本地开发（一键启动）

```bash
# 1. 克隆项目
git clone <your-repo-url>
cd Learn4Go-1

# 2. 一键启动（内存模式）
./start-local.sh

# 3. 访问服务
# 学习门户: http://localhost:8000/portal.html
# TODO API: http://localhost:8080
# Gateway: http://localhost:8888

# 4. 停止服务
./stop-local.sh

# 5. 清理临时文件
./clean.sh
```

**其他启动模式**：

```bash
# SQLite 模式
./start-local.sh sqlite

# MySQL 模式（需要先启动 Docker MySQL）
docker-compose -f deployments/docker-compose.yml up -d mysql
./start-local.sh mysql
```

详细说明请参考 [本地开发指南](docs/本地开发指南.md)。

## 📂 项目结构

```
Learn4Go-1/
├── 01.Go语言基础/              # 10 个核心章节
│   ├── 01_快速开始与基本语法.md
│   ├── 02_变量_常量_类型.md
│   └── ...
├── 02.开发环境及框架介绍/       # 框架与工具链
│   ├── 01_工具链_依赖管理.md
│   ├── 02_Web框架对比_Gin_Echo_Fiber.md
│   ├── 03_gRPC_与微服务入门.md
│   ├── 04_日志_配置_中间件.md
│   ├── 05_配置与日志.md         # Viper + Zap
│   └── 06_数据库访问.md         # GORM
├── 03.项目实战/                # 实战项目说明
├── cmd/                        # 可执行程序入口
│   ├── todoapi/               # TODO REST API
│   ├── batchrename/           # 批量重命名 CLI
│   └── logzip/                # 日志归档工具
├── internal/                   # 内部模块
│   ├── todo/                  # TODO 业务逻辑
│   │   ├── handler.go         # HTTP 处理器
│   │   ├── auth.go            # JWT 认证
│   │   ├── user_store.go      # 用户存储
│   │   ├── ratelimit.go       # 速率限制
│   │   ├── store.go           # 内存存储
│   │   └── store_db.go        # 数据库存储
│   ├── cache/                 # Redis 缓存
│   │   └── redis.go
│   └── storage/               # MinIO 对象存储
│       └── minio.go
├── examples/                   # 示例代码
│   ├── hello/                 # Hello World
│   ├── workerpool/            # 并发任务池
│   ├── grpc/                  # gRPC 示例
│   │   ├── server/
│   │   └── client/
│   ├── gateway/               # API 网关
│   │   ├── stdlib/            # 标准库实现
│   │   └── gin/               # Gin 实现
│   ├── basics/                # 语言基础练习（指针/切片/map/defer/error/channel/context 等）
│   ├── java_compare/          # Java→Go 差异练习（并发/接口/超时/urlencode/httptrace/词频/ticker/中间件/pprof/sync.Map）
│   └── tinygee/               # 手撸微框架 TinyGee 的分日示例
│   ├── config/                # Viper 配置示例
│   └── database/              # GORM 数据库示例
├── web/                        # 前端页面
│   ├── portal.html            # 学习门户
│   ├── config.js              # 前端配置
│   └── ...
├── deployments/                # 部署配置
│   ├── docker-compose.yml     # Docker 编排
│   ├── Dockerfile             # TODO API 镜像
│   ├── Dockerfile.gateway     # Gateway 镜像
│   ├── nginx.conf             # Nginx 配置
│   ├── README.md              # 部署文档
│   └── .env.example           # 环境变量示例
├── docs/                       # 文档
│   ├── Java_vs_Go_CheatSheet.md  # Java/Go 对照表
│   └── ...
└── plan.md                     # 项目计划
```

## 🎯 学习路径

### 阶段一：Go 语言基础（1-2 周）

学习 `01.Go语言基础/` 目录下的 10 个章节：

1. 快速开始与基本语法
2. 变量、常量、类型
3. 流程控制
4. 函数与错误处理
5. 数组、切片、Map
6. 结构体、方法、组合
7. 接口与多态
8. 并发：goroutine 与 channel
9. 包与模块管理
10. 测试与基准

**实践项目**：

- `examples/hello/` - Hello World
- `examples/fizzbuzz/` - FizzBuzz 练习
- `cmd/batchrename/` - 批量重命名 CLI
- `examples/java_compare/` - Java→Go 对照练习（并发/接口/超时/urlencode/httptrace/词频/ticker/中间件/pprof/sync.Map）
- `examples/basics/` - 语言基础练习（指针接收者、切片底层、map 零值、defer 捕获、错误包装、channel 超时、context 取消）
- `examples/io/` - IO/压缩/JSON/XML 示例
- `examples/network/` - TCP/UDP/http client+ctx 示例
- `examples/testing/benchmark_pprof/` - 基准与 pprof 入门
- `examples/tinygee/` - 手撸框架 TinyGee 分日示例
- `tinygee/` + `examples/tinygee/` - 参考 Gee 教程的最小 Web 框架实现与分日示例（路由、中间件、模板、JWT/RBAC、限流、Recover）

更多练习清单见 `docs/Java对照练习.md`。

## 基础学习资源

- 推荐阅读 `docs/Go学习资源.md`：包含官方文档入口、中文教程（含菜鸟教程对照）、学习顺序与仓库内配套练习指引。搭配 `docs/Go基础速查.md` 快速回顾要点与对应示例。

## TinyGee 快速上手

- 目标：在本仓库中手撸一个最小 Web 框架（无第三方依赖），并演示路由、中间件、模板、JWT/RBAC、限流等能力。
- 文档：`docs/tinygee/quickstart.md`、`docs/tinygee/tinygee_execution_plan.md`
- 运行 demo：`go run ./cmd/tinygee-demo --port :9999 --prom=true`
- 查看分日示例：`go run ./examples/tinygee/day1`（基础路由）、`day3`（动态路由）、`day5`（分组中间件）、`day6`（模板/静态）、`day7`（Recover）、`day9auth`（JWT/RBAC）、`day9rl`（限流）

### 阶段二：框架与工具链（1 周）

学习 `02.开发环境及框架介绍/` 目录：

1. 工具链与依赖管理（go mod）
2. Web 框架对比（Gin vs Echo vs Fiber）
3. gRPC 与微服务入门
4. 日志、配置与中间件
5. Viper 配置管理
6. GORM 数据库访问

**实践项目**：

- `examples/config/` - Viper 配置示例
- `examples/database/` - GORM 数据库示例
- `examples/grpc/` - gRPC 客户端/服务端

### 阶段三：微服务实战（2-3 周）

构建完整的微服务应用：

1. **TODO REST API** (`cmd/todoapi/`)

   - 支持内存/SQLite/MySQL 存储
   - JWT 认证与授权 🔐
   - 健康检查端点
   - 环境变量配置

2. **认证系统** (`internal/todo/auth.go`)

   - 用户注册/登录
   - JWT token 生成与验证
   - bcrypt 密码加密
   - 速率限制

3. **API 网关** (`examples/gateway/gin/`)

   - 反向代理
   - 请求日志
   - 限流与熔断（可扩展）

4. **缓存层** (`internal/cache/`)

   - Redis 集成
   - 限流实现

5. **对象存储** (`internal/storage/`)
   - MinIO 集成
   - 文件上传/下载

### 阶段四：部署与运维（1 周）

学习 Docker 容器化部署：

1. 阅读 `deployments/README.md`
2. 理解 Docker Compose 编排
3. 学习健康检查与重启策略
4. 实践日志收集与监控

## 📖 核心文档

| 文档                                              | 说明                |
| ------------------------------------------------- | ------------------- |
| [常见问题 FAQ](docs/常见问题.md)                  | 常见问题解答 ⭐     |
| [Java vs Go CheatSheet](docs/Java对照Go速查表.md) | Java 开发者速查表   |
| [本地开发指南](docs/本地开发指南.md)              | 本地启动详细步骤    |
| [部署指南](deployments/README.md)                 | Docker 部署完整文档 |
| [API 文档](docs/API接口文档.md)                   | TODO API 接口文档   |
| [JWT 认证文档](docs/JWT认证系统.md)               | JWT 认证系统说明 🔐 |
| [变更日志](docs/变更日志.md)                      | 版本更新记录        |
| [前端使用指南](docs/前端使用指南.md)              | 前端页面使用说明    |
| [项目计划](plan.md)                               | 后续开发计划        |

## 🏗️ 技术栈

### 后端

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **RPC 框架**: gRPC
- **ORM**: GORM
- **配置管理**: Viper
- **日志**: 标准库 log（可扩展 Zap）
- **认证**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt (golang.org/x/crypto)
- **速率限制**: 内存滑动窗口

### 基础设施

- **数据库**: MySQL 8
- **缓存**: Redis 7
- **对象存储**: MinIO
- **容器**: Docker + Docker Compose
- **反向代理**: Nginx

### 前端

- **框架**: 原生 HTML/CSS/JavaScript
- **UI**: Bootstrap 5
- **字体**: Inter + Fira Code

## 🎨 前端功能

访问 http://localhost 查看学习门户：

- **学习进度追踪**: 自动保存学习进度到浏览器
- **章节导航**: 基础/框架/项目三大模块
- **服务状态监控**: 实时显示后端服务状态
- **代码示例**: 可复制运行的示例命令
- **响应式设计**: 支持桌面和移动端

## 🔧 开发工具

```bash
# 代码格式化
go fmt ./...

# 静态检查
go vet ./...

# 运行测试
go test ./...

# 构建所有可执行文件
go build ./cmd/...

# 清理构建缓存
go clean -cache
```

## 📊 项目状态

- ✅ Go 语言基础教程（10 章节）
- ✅ 框架介绍文档（6 章节）
- ✅ TODO REST API（内存/数据库双模式）
- ✅ API 网关（标准库/Gin 双实现）
- ✅ Redis 缓存集成
- ✅ MinIO 对象存储集成
- ✅ Docker Compose 部署
- ✅ 前端学习门户
- ✅ Java vs Go 速查表

## 🚧 后续计划

详见 [plan.md](plan.md)：

- [ ] 添加 JWT 认证中间件
- [ ] 集成 Prometheus 监控
- [ ] 添加 Grafana 可视化
- [ ] 实现分布式追踪（Jaeger）
- [ ] 添加单元测试覆盖
- [ ] 性能压测与优化
- [ ] Kubernetes 部署示例

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- Go 官方文档
- Gin Web Framework
- GORM ORM
- Docker 社区

## 📮 联系方式

如有问题或建议，欢迎通过 Issue 反馈。

---

**Happy Coding! 🎉**
