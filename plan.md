# Learn4Go 迭代计划

> 面向 Java 开发者的 Go 语言系统化学习路线

## 2025-12-04 状态快照与路线图

- 🎯 项目目的：为 Java 开发者提供从 Go 语言基础到微服务实践、部署运维的一站式学习路径，涵盖代码示例、前端门户、Docker 一键编排。
- 🗂️ 目录/模块速览：`01.Go语言基础`、`02.开发环境及框架介绍`、`03.项目实战`、`cmd/*` 可执行程序、`examples/*` 示例、`internal/*` 业务实现、`deployments/` 容器化、`web/` 学习门户。

### 当前完成度
- ✅ 基础文档与速查表：顶层 README、Java_vs_Go_CheatSheet、部署/前端指南
- ✅ 示例与实战：CLI 批量重命名、日志归档、TODO API（内存+GORM）、gRPC（unary/streaming）、Gateway（stdlib+Gin）
- ✅ 基础设施：Docker Compose、Redis/MinIO 集成、Makefile（fmt/vet/test）
- ⏳ 工程增强：golangci-lint、GitHub Actions CI、单元测试覆盖率 >80%
- ⏳ 架构升级：依赖注入（Wire）、JWT/RBAC、Prometheus 指标、OpenTelemetry 追踪
- ⏳ 运维与安全：压测基准、限流/熔断、HTTPS/TLS、Swagger/OpenAPI 规范

### 接下来 1-2 周优先级（建议）
1. 工程质量：补充 `golangci-lint` 配置 + GitHub Actions（go fmt/vet/test + lint）。
2. 稳定性：为 `cmd/todoapi` 与网关添加 httptest 覆盖，重点覆盖存储切换与错误路径。
3. 安全基线：实现 JWT 登录 + 路由鉴权中间件，并在 Gateway 演示限流/熔断基础版。
4. 可观测性：添加基础 Prometheus 指标（HTTP 请求计数/延迟、DB 连接池），预留 OpenTelemetry 接口。
5. 文档一致性：为新增能力补充 README/FRONTEND/API 说明，保持学习路径闭环。

### 开发技能点（勾选=已覆盖）
- 语言/基础：✅ 并发、接口、多态、测试基准
- 框架：✅ Gin / gRPC；⏳ Wire DI；⏳ gRPC 拦截器链
- 数据与存储：✅ GORM、Redis、MinIO；⏳ 消息队列（Kafka/RabbitMQ）
- 工程化：✅ fmt/vet/test；⏳ golangci-lint；⏳ CI/CD（GitHub Actions）
- 观测与性能：⏳ Prometheus/Grafana；⏳ OpenTelemetry + Jaeger；⏳ k6/wrk 压测
- 安全：⏳ JWT/RBAC；⏳ HTTPS/TLS；⏳ 限流/熔断；⏳ Swagger/OpenAPI


## 项目结构说明

```
Learn4Go-1/
├── 01.Go语言基础/          # 10 章语法与并发基础文档
├── 02.开发环境及框架介绍/   # 工具链、Web框架、gRPC、中间件
├── 03.项目实战/            # 实战项目说明文档
├── cmd/                    # 可执行入口 (Go 标准布局)
│   ├── todoapi/           #   → 内存版 TODO REST API
│   ├── batchrename/       #   → 批量重命名 CLI 工具
│   ├── logzip/            #   → 日志 ZIP 归档工具
│   └── learn4go/          #   → 最简入口示例
├── internal/              # 业务逻辑 (不可被外部导入)
│   ├── todo/              #   → TODO API handler/store
│   ├── workerpool/        #   → 并发 Worker 池
│   ├── cli/batchrename/   #   → 重命名核心逻辑
│   └── logzip/            #   → 日志压缩逻辑
├── examples/              # 可运行示例代码
│   ├── variables/         #   → 变量/常量/类型
│   ├── controlflow/       #   → 流程控制
│   ├── functions/         #   → 函数/闭包/defer
│   ├── collections/       #   → 数组/切片/map
│   ├── structs/           #   → 结构体/方法
│   ├── interfaces/        #   → 接口/多态
│   ├── concurrency/       #   → goroutine/channel
│   ├── packages/          #   → 标准库使用
│   ├── testing/           #   → 单元测试/基准测试
│   ├── workerpool/        #   → Worker 池使用示例
│   ├── grpc/              #   → gRPC unary + streaming
│   └── gateway/           #   → API 网关 (stdlib + Gin)
├── web/                   # 前端页面
│   └── portal.html        #   → 学习门户页面
└── docs/                  # 文档
    ├── Go学习规划_Java开发者版.md
    ├── Java_vs_Go_CheatSheet.md
    └── architecture.md
```

## 当前进度

### ✅ 已完成（Phase 1 - 基础架构）
- ✅ 仓库 scaffold：基础目录、README、MIT License
- ✅ Go 标准布局：`cmd/*` 入口、`internal/*` 业务逻辑
- ✅ 基础章节示例：9 个 examples (variables → testing)
- ✅ 并发 Worker 池：`internal/workerpool` + 测试
- ✅ gRPC 示例：unary RPC + server streaming
- ✅ API 网关示例：标准库版 + Gin 版
- ✅ 前端门户：Bootstrap 5 响应式页面
- ✅ 实战示例：batchrename、todoapi、logzip
- ✅ Java vs Go 速查表（10 大核心对比）
- ✅ Viper 配置管理文档与示例
- ✅ GORM 数据库访问文档与示例
- ✅ TODO API 数据库版本（MySQL/SQLite）
- ✅ Redis 缓存集成（`internal/cache/redis.go`）
- ✅ MinIO 对象存储集成（`internal/storage/minio.go`）
- ✅ Docker Compose 微服务编排
- ✅ 完整部署文档（`deployments/README.md`）
- ✅ API 接口文档（`docs/API.md`）
- ✅ 前端使用指南（`docs/FRONTEND.md`）
- ✅ 错误处理优化（存储接口返回 error）

### ⏳ 进行中（Phase 2 - 功能增强）
- 🔄 单元测试覆盖率提升
- 🔄 性能压测与优化

---

## Java → Go 学习路线图

### 阶段 1：语言基础 (1-2 周)
**目标**：建立 Go 思维，理解与 Java 的核心差异

| 章节 | 内容 | Java 对照 | 示例 |
|------|------|-----------|------|
| 01 | 快速开始 | public static void main | `examples/hello` |
| 02 | 变量/常量/类型 | 基本类型、类型推断 | `examples/variables` |
| 03 | 流程控制 | if/switch/for | `examples/controlflow` |
| 04 | 函数与错误处理 | 方法、try-catch → if err | `examples/functions` |
| 05 | 数组/切片/map | ArrayList/HashMap | `examples/collections` |
| 06 | 结构体与方法 | class/extends → struct/embedding | `examples/structs` |
| 07 | 接口与多态 | interface implements → 隐式实现 | `examples/interfaces` |
| 08 | 并发 | ThreadPool → goroutine+channel | `examples/concurrency` |
| 09 | 包与模块 | Maven/Gradle → go mod | `examples/packages` |
| 10 | 测试与基准 | JUnit → go test | `examples/testing` |

**Java 开发者常见坑点**：
- `nil` map 写入会 panic，需显式 `make()`
- 切片 `append` 可能重新分配底层数组
- 值接收者 vs 指针接收者的选择
- `defer` 在循环中的开销
- 无异常机制，错误需显式处理

### 阶段 2：框架与工程化 (1-2 周)
**目标**：掌握 Go 生态主流框架，对标 Spring Boot

#### 框架学习顺序（推荐）

```
1. Gin (HTTP)          ← 类比 Spring MVC
   ↓
2. Viper (配置)        ← 类比 Spring Config
   ↓
3. Zap/slog (日志)     ← 类比 Logback/SLF4J
   ↓
4. GORM/sqlc (数据库)  ← 类比 JPA/MyBatis
   ↓
5. Wire (依赖注入)     ← 类比 Spring DI
   ↓
6. gRPC (微服务)       ← 类比 OpenFeign
   ↓
7. 中间件 (横切关注点)  ← 类比 Spring AOP
```

#### 框架详细学习步骤

**Step 1: Gin Web 框架**
```bash
# 1. 安装
go get github.com/gin-gonic/gin

# 2. 学习内容
- 路由定义与分组
- 请求参数绑定 (c.ShouldBindJSON)
- 中间件链 (Logger, Recovery, CORS)
- 响应处理 (c.JSON, c.String)

# 3. 练习
- 将 cmd/todoapi 从 net/http 迁移到 Gin
- 添加请求日志中间件
```

**Step 2: Viper 配置管理**
```bash
go get github.com/spf13/viper

# 学习内容
- 读取 YAML/JSON/ENV 配置
- 配置热重载
- 环境变量覆盖
```

**Step 3: 结构化日志**
```bash
go get go.uber.org/zap
# 或使用 Go 1.21+ 标准库 slog

# 学习内容
- 结构化字段日志
- 日志级别控制
- 与 Gin 集成
```

**Step 4: 数据库访问**
```bash
# 方案 A: GORM (类似 JPA)
go get gorm.io/gorm
go get gorm.io/driver/sqlite

# 方案 B: sqlc (类似 MyBatis)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# 学习内容
- 模型定义与迁移
- CRUD 操作
- 事务处理
- 连接池配置
```

**Step 5: 依赖注入**
```bash
go install github.com/google/wire/cmd/wire@latest

# 学习内容
- Provider 定义
- Injector 生成
- 分层架构 (handler/service/repo)
```

**Step 6: gRPC 微服务**
```bash
# 已有示例: examples/grpc/
go get google.golang.org/grpc
go get google.golang.org/protobuf

# 学习内容
- Proto 文件定义
- Unary RPC vs Streaming
- 拦截器 (类似 Spring AOP)
- 错误码处理
```

### 阶段 3：项目实战 (2-3 周)
**目标**：通过渐进式项目巩固知识

| 序号 | 项目 | 技术点 | 入口 |
|------|------|--------|------|
| 1 | CLI 批量重命名 | flag/cobra、文件操作、测试 | `cmd/batchrename` |
| 2 | 日志归档工具 | archive/zip、crypto、IO | `cmd/logzip` |
| 3 | TODO REST API | net/http → Gin、内存存储 | `cmd/todoapi` |
| 4 | TODO + 数据库 | GORM/sqlc、迁移、事务 | 待实现 |
| 5 | Worker 任务队列 | goroutine、channel、超时 | `examples/workerpool` |
| 6 | gRPC TODO 服务 | protobuf、streaming、拦截器 | `examples/grpc` |
| 7 | API 网关 | 反向代理、JWT、限流 | `examples/gateway` |

### 阶段 4：进阶与运维
- Prometheus 指标埋点
- OpenTelemetry 链路追踪
- Docker 容器化部署
- GitHub Actions CI/CD

---

## Java → Go 核心概念映射

| Java | Go | 说明 |
|------|-----|------|
| `class` | `struct` | Go 无类继承 |
| `extends` | 嵌入 (embedding) | 组合优于继承 |
| `implements` | 隐式实现 | 无需声明 |
| `@Autowired` | Wire/Fx | 静态/运行时 DI |
| `try-catch` | `if err != nil` | 显式错误处理 |
| `ThreadPoolExecutor` | goroutine + channel | CSP 并发模型 |
| `synchronized` | `sync.Mutex` | 互斥锁 |
| `Future/Promise` | channel + select | 异步通信 |
| `@RestController` | Gin HandlerFunc | HTTP 处理器 |
| `@RequestBody` | `c.ShouldBindJSON` | 请求绑定 |
| `JPA Entity` | GORM Model | ORM 映射 |
| `@Transactional` | `db.Transaction()` | 事务处理 |
| `Logback` | Zap/slog | 结构化日志 |
| `application.yml` | Viper | 配置管理 |

---

## 任务清单

### 文档
- [x] 基础目录结构文档
- [x] 完善 `docs/Java_vs_Go_CheatSheet.md` 代码对照
- [x] API 接口文档 `docs/API.md`
- [x] 前端使用指南 `docs/FRONTEND.md`
- [x] Docker 部署文档 `deployments/README.md`
- [ ] 补充各章节练习题

### 框架章节
- [x] `02.开发环境及框架介绍/05_配置与日志.md` (Viper + Zap)
- [x] `02.开发环境及框架介绍/06_数据库访问.md` (GORM)
- [x] `examples/config/` - Viper 配置示例
- [x] `examples/database/` - GORM 数据库示例
- [ ] `02.开发环境及框架介绍/07_依赖注入.md` (Wire)

### 实战项目
- [x] CLI 批量重命名
- [x] 日志归档工具
- [x] TODO REST API (内存版)
- [x] TODO + MySQL/SQLite (GORM 版)
- [x] Worker 池示例
- [x] gRPC 示例
- [x] API 网关示例
- [x] 完整微服务示例 (TODO + gRPC + 网关 + Redis + MinIO)
- [x] Docker Compose 微服务编排

### 质量
- [x] Makefile (fmt/vet/test)
- [x] 错误处理优化（存储接口返回 error）
- [ ] 单元测试覆盖率 >80%
- [ ] golangci-lint 配置
- [ ] GitHub Actions CI

---

## 运行指南

```bash
# 基础示例
go run ./examples/variables
go run ./examples/concurrency

# 实战项目
go run ./cmd/todoapi              # TODO API :8080
go run ./cmd/batchrename --help   # CLI 工具
go run ./cmd/logzip               # 日志归档

# 网关示例
go run ./examples/gateway/gin     # Gin 网关 :8888
go run ./examples/gateway/stdlib  # 标准库网关 :8888

# 测试
go test ./...
go test -bench=. ./examples/testing

# 代码质量
make fmt vet test
```

## 后续计划（Phase 3 - 生产就绪）

### 🎯 短期目标（1-2 周）

#### 1. 认证与授权
- [x] JWT 中间件实现 ✅ (v1.1.0)
- [x] 用户注册/登录 API ✅ (v1.1.0)
- [x] bcrypt 密码加密 ✅ (v1.1.0)
- [x] 速率限制（内存版）✅ (v1.1.0)
- [ ] 基于角色的访问控制（RBAC）
- [ ] 前端登录页面集成
- [ ] Token 刷新机制
- [ ] 登录失败次数限制

#### 2. 测试覆盖
- [ ] TODO API 单元测试（目标 >80%）
- [ ] Gateway 单元测试
- [ ] Redis/MinIO 集成测试
- [ ] E2E 测试（使用 httptest）

#### 3. 监控与可观测性
- [ ] Prometheus 指标埋点
  - HTTP 请求计数/延迟
  - 数据库连接池状态
  - Redis 命中率
- [ ] Grafana 仪表盘
- [ ] 健康检查增强（依赖检查）

#### 4. 性能优化
- [ ] 压测基准（wrk/k6）
- [ ] 数据库查询优化
- [ ] Redis 缓存策略
- [ ] 连接池调优

### 🚀 中期目标（3-4 周）

#### 5. 分布式追踪
- [ ] OpenTelemetry 集成
- [ ] Jaeger 部署
- [ ] 跨服务链路追踪

#### 6. 消息队列
- [ ] RabbitMQ/Kafka 集成
- [ ] 异步任务处理
- [ ] 事件驱动架构示例

#### 7. 服务网格
- [ ] Istio 入门示例
- [ ] 服务间通信加密
- [ ] 流量管理与灰度发布

#### 8. CI/CD
- [ ] GitHub Actions 工作流
  - 自动测试
  - 代码质量检查（golangci-lint）
  - Docker 镜像构建
  - 自动部署到测试环境

### 🌟 长期目标（1-2 个月）

#### 9. Kubernetes 部署
- [ ] Helm Chart 编写
- [ ] K8s 资源定义（Deployment/Service/Ingress）
- [ ] ConfigMap/Secret 管理
- [ ] HPA 自动扩缩容
- [ ] 滚动更新与回滚

#### 10. 高级特性
- [ ] 分布式锁（Redis/etcd）
- [ ] 分布式事务（Saga 模式）
- [ ] 限流与熔断（Sentinel）
- [ ] API 版本管理
- [ ] GraphQL 支持

#### 11. 安全加固
- [ ] HTTPS/TLS 配置
- [ ] API 限流（令牌桶/漏桶）
- [ ] SQL 注入防护审计
- [ ] XSS/CSRF 防护
- [ ] 敏感数据加密

#### 12. 文档完善
- [ ] Swagger/OpenAPI 规范
- [ ] 架构决策记录（ADR）
- [ ] 运维手册
- [ ] 故障排查指南

### 📚 学习资源扩展

#### 13. 进阶教程
- [ ] Go 并发模式深入
- [ ] Go 内存模型与 GC
- [ ] Go 性能分析（pprof）
- [ ] Go 汇编与底层优化

#### 14. 实战案例
- [ ] 短链接服务
- [ ] 实时聊天系统（WebSocket）
- [ ] 文件上传/下载服务
- [ ] 定时任务调度系统

---

## 🎓 学习建议

### 对于初学者
1. 先完成 `01.Go语言基础/` 的 10 个章节
2. 运行 `examples/` 下的所有示例代码
3. 阅读 `docs/Java_vs_Go_CheatSheet.md` 对照学习
4. 完成 `cmd/batchrename` 和 `cmd/logzip` 小项目

### 对于有经验的开发者
1. 直接从 `02.开发环境及框架介绍/` 开始
2. 研究 `cmd/todoapi` 的完整实现
3. 部署 Docker Compose 环境，理解微服务架构
4. 参与后续计划的功能开发

### 对于架构师
1. 研究 `deployments/` 下的部署配置
2. 理解微服务间的通信模式
3. 规划 Kubernetes 迁移方案
4. 设计监控与告警体系

---

## 📞 贡献与反馈

欢迎通过以下方式参与项目：

1. **提交 Issue**: 报告 Bug 或提出新功能建议
2. **Pull Request**: 贡献代码或文档
3. **讨论**: 在 Discussions 中分享学习心得
4. **Star**: 如果觉得有帮助，请给项目点个 Star ⭐

---

**最后更新**: 2025-12-05
