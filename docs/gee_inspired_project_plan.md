# 基于 Gee 教程的全新实战计划（草案）

## 背景
- 参考极客兔兔《7 天用 Go 从零实现 Web 框架 Gee》系列，结合当前仓库的 TODO/Gateway/认证实践，设计一条“从零手撸微框架 → 对标生产特性”的学习/实战线。
- 目标：在不依赖第三方 Web 框架的前提下，用 1~2 周迭代产出一个可运行的极简框架（暂命名 **TinyGee**），并逐步嫁接现有组件（JWT/RBAC/中间件/模板/日志/链路追踪）。

## 参考来源
- Gee 系列目录：Day1 http.Handler、Day2 Context、Day3 Trie 路由、Day4 分组、Day5 中间件、Day6 模板、Day7 Panic Recover。citeturn0search0turn0search2turn0search3turn0search6

## 目录隔离约定
- 所有 TinyGee 代码与示例均置于 `tinygee/` 顶级目录下的子包与 `examples/tinygee/`，不与现有 `cmd/`、`internal/`、`examples/`（其他子目录）互相引用。  
- 新可执行入口放在 `cmd/tinygee-demo/`，仅依赖 `tinygee/` 包和标准库；不引入 gin。  
- 文档放在 `docs/tinygee/`，与现有文档平行。

## 总体里程碑（建议 10 天）
1. **基础内核 (Day1-2)**  
   - 目标：可处理请求、包装 Context、返回 JSON/String。  
   - 交付：`tinygee/context.go`、`tinygee/engine.go`、最小示例 `examples/tinygee/day1_main.go`。
2. **路由与动态匹配 (Day3)**  
   - 目标：Trie 路由，支持 `:param`、`*filepath`；路由查找性能基准。  
   - 交付：`tinygee/router.go` + bench；示例覆盖动态/通配。
3. **路由分组与中间件 (Day4-5)**  
   - 目标：Group 前缀、分组级中间件、全局 Logger/Recovery。  
   - 交付：`tinygee/group.go`、`tinygee/middleware.go`；示例 `logger` / `auth stub`。
4. **模板与静态资源 (Day6)**  
   - 目标：`SetFuncMap`、`LoadHTMLGlob`、静态文件服务。  
   - 交付：`tinygee/render.go`；示例模板加载、静态目录。
5. **错误恢复与健壮性 (Day7)**  
   - 目标：Panic Recover 中间件、统一错误响应、可配置。  
   - 交付：`tinygee/recover.go`；故障注入示例。
6. **扩展对接现有能力 (Day8-10)**  
   - 集成已有 JWT/RBAC：将 `internal/todo` 的 JWT/权限逻辑“复制改写”成 TinyGee middleware（新文件放 `tinygee/middleware/auth`），不直接 import 原代码，保持隔离。  
   - 集成限流/熔断：同理，将限流逻辑迁移为 `tinygee/middleware/ratelimit`，避免跨依赖。  
   - 可观测性：最小 Prometheus 指标（请求计数/延迟）接口挂载。  
   - 配置与启动：提供 `cmd/tinygee-demo`，可通过 flag/env 选择端口、模板路径。  
   - 测试：`httptest` 覆盖动态路由、分组中间件、Recover、模板渲染、权限拒绝。

## 任务拆解（执行版，完全隔离）
- 代码结构：  
  - 核心包：`tinygee/engine.go`、`context.go`、`router.go`、`group.go`、`middleware/logger.go`、`middleware/recover.go` 等  
  - 扩展包：`tinygee/middleware/auth`、`tinygee/middleware/ratelimit`、`tinygee/render`（模板）  
  - 可执行入口：`cmd/tinygee-demo/main.go`  
  - 示例：`examples/tinygee/day1_basic.go`、`day3_routes.go`、`day5_middlewares.go`、`day7_recover.go`、`day9_auth.go` 等  
  - 文档：`docs/tinygee/` 下 quickstart/design/middleware/auth/render/bench
- Day1-2：  
  - [ ] `Engine` 实现 `http.Handler`；`Context` 封装 Req/Res/Params；JSON/String/Status 辅助。  
  - [ ] 示例与基准：Hello、Query/POST 解析。
- Day3：  
  - [ ] `node` + Trie 插入/匹配；解析 `:param` `*filepath`。  
  - [ ] 基准：1000 条路由匹配耗时；对比 `map[string]Handler`。
- Day4-5：  
  - [ ] `RouterGroup` 前缀叠加；`Use` 注册中间件；`Context.handlers + Next()` 链式。  
  - [ ] 默认中间件：Logger（耗时/状态码）、Recovery（stack）。  
  - [ ] 示例：全局 Logger + v1/v2 组限流/鉴权。
- Day6：  
  - [ ] `SetFuncMap`、`LoadHTMLGlob`、`Static`；加入简单模板示例。  
- Day7：  
  - [ ] Panic Recover 覆盖所有 handler；可配置响应格式（JSON/HTML）。  
- Day8-10（增强与整合）：  
  - [ ] JWT/RBAC 适配：将 `internal/todo` 中 `authMiddleware`、`authzMiddleware` 改写为 TinyGee middleware。  
  - [ ] RateLimiter/熔断：把 `internal/todo` 的限流包装，增加按 IP/路由维度 key 选择。  
  - [ ] Prometheus：暴露 `/metrics`（可选 `promhttp.Handler()` 集成）。  
  - [ ] CLI：`cmd/tinygee-demo` 支持 `--port --static --templates --prom`。  
  - [ ] 文档：新增 `docs/tinygee/*.md`，说明路由、上下文、中间件、模板、恢复、鉴权。  
  - [ ] 测试：`go test ./tinygee/...` + `httptest`；性能基准 `go test -bench=.`。

## 与现有仓库的结合
- 复用与迁移：  
  - `internal/todo` 的 JWT/RBAC/限流作为中间件示例，展示“框架无关的横切逻辑如何插拔”。  
  - `docs/API.md` 可新增一节“TinyGee + TODO”对接示例。  
- 前端：保留 `todo-login.html`，新增一条使用 TinyGee demo API 的配置项（后续实现时补）。
- CI：新增 `make tinygee-test` 执行 `go test ./tinygee/...`；增加基准命令可选。

## 工具与规范
- 依赖：仅使用标准库；增强阶段可选 `prometheus/client_golang`。  
- 测试：表格驱动 + `httptest`；基准使用 `testing.B`。  
- 代码质量：保持中文注释、gofmt、避免引入非必要依赖。

## 可交付清单（落地时）
- 代码：`tinygee/` 源码、`cmd/tinygee-demo` 入口、`examples/tinygee/*`。  
- 测试：单测 + 基准。  
- 文档：`docs/tinygee/quickstart.md`、`docs/tinygee/middleware.md`、`docs/tinygee/design.md`、与现有 `README/FRONTEND/API` 的对接说明。  
- 说明：`docs/go_basics_cheatsheet.md`、`docs/go_basics_resources.md` 可作为前置阅读链接。
