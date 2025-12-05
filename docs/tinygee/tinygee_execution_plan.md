# TinyGee 执行计划（隔离版）

## 目录与依赖约束
- 代码仅放在 `tinygee/` 子目录；示例放 `examples/tinygee/`；入口放 `cmd/tinygee-demo/`。  
- 不引用现有 `internal/`、`cmd/`、`examples/` 下的实现；需要的功能（JWT、限流）以“复制+改写”方式独立实现。  
- 依赖：标准库为主；可选 `prometheus/client_golang`（若引入则仅 TinyGee 使用）。  

## 里程碑与交付
### Day1-2 基础内核
- 目标：`Engine` 实现 `http.Handler`；`Context` 包装请求/响应；JSON/String/Status 辅助。
- 交付：✅ `tinygee/engine.go`、`context.go`、`router.go`；示例 `examples/tinygee/day1_basic.go`。
- 测试：`go test ./tinygee/...` 已通过。

### Day3 路由 Trie
- 目标：支持 `:param`、`*filepath` 动态路由；路由查找基准。
- 交付：✅ `tinygee/trie.go`、`router_test.go`、示例 `examples/tinygee/day3_routes.go`。
- 测试：表格驱动已覆盖参数与通配符；基准待后续 `tinygee-bench`。

### Day4-5 分组与中间件
- 目标：`RouterGroup` 前缀叠加；中间件链（logger、recover）；`Context.Next()`。
- 交付：✅ `tinygee/group.go`、`middleware/logger.go`、`middleware/recover.go`；示例 `examples/tinygee/day5_middlewares.go`。
- 测试：路由 & 中间件顺序在 `router_test.go` 覆盖，panic 由 recover 中间件处理。

### Day6 模板与静态
- 目标：`SetFuncMap`、`LoadHTMLGlob`、`Static`。
- 交付：✅ `tinygee/render/render.go`；示例 `examples/tinygee/day6_template.go` + 模板/静态资源。
- 测试：后续可补 httptest 覆盖模板/静态。

### Day7 健壮性
- 目标：统一错误响应、Recover 配置（JSON/HTML）。
- 交付：✅ `middleware/recover.go` 支持 JSON/文本模式；示例 `examples/tinygee/day7_recover.go`。
- 测试：待补 HTML/自定义模式（当前涵盖 JSON）。

### Day8-9 安全与限流（隔离实现）
- 目标：JWT + RBAC 中间件；RateLimiter 中间件。
- 交付：✅ `middleware/auth/jwt.go`、`middleware/auth/rbac.go`、`middleware/ratelimit/memory.go`；示例 `day9_auth.go`、`day9_ratelimit.go`。
- 测试：✅ JWT 过期/合法用例；✅ 限流拒绝；RBAC 角色拒绝待补。

### Day10 可观测与示例整合
- 目标：可选 Prometheus `/metrics`；CLI 入口。
- 交付：✅ `cmd/tinygee-demo` 支持 `--port` 与可选 `/metrics`（expvar）；`examples/tinygee/fullstack` 已整合模板+JWT/RBAC+限流+metrics；`engine_metrics_test.go` 覆盖 /metrics。
- 测试：`go test ./tinygee/...` 覆盖 metrics；fullstack 为运行示例。

## 文档与工具
- 文档目录：`docs/tinygee/quickstart.md`、`design.md`、`middleware.md`、`auth.md`、`benchmarks.md`。
- Makefile：新增 `tinygee-test`、`tinygee-bench`（不影响现有目标）。
- CI：后续可在 GitHub Actions 中添加独立 job（与现有分离）。

## 执行顺序（可开始动手）
1) 创建目录骨架：`tinygee/`、`tinygee/middleware`、`tinygee/render`、`cmd/tinygee-demo/`、`examples/tinygee/`、`docs/tinygee/`。  
2) 实现 Day1-2 内核 + 基础测试。  
3) 迭代 Day3 路由 + 测试 + bench。  
4) Day4-5 中间件/分组；Day6 模板；Day7 Recover 增强。  
5) Day8-9 JWT/RBAC/限流（独立实现）。  
6) Day10 Prom + demo 入口 + 文档补全。  
7) go test ./tinygee/...；更新 README/plan 与门户示例（可选）。  

## 注意
- 全程避免跨目录依赖，保持 TinyGee 可单独抽离。  
- 示例使用标准库和 TinyGee，不引入 gin。  
- 若需对比 gin，可放在 `examples/tinygee/compare_gin.go`，但确保 import gin 仅在该文件/示例，不进入 TinyGee 核心。
