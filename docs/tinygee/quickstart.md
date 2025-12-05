# TinyGee 快速开始

## 目录
- 代码：`tinygee/`（核心）、`tinygee/middleware/`（logger/recover/auth/ratelimit）、`tinygee/render/`（模板）、`cmd/tinygee-demo/`（入口）、`examples/tinygee/`（示例）
- 文档：`docs/tinygee/`（本文件）、`tinygee_execution_plan.md`

## TinyGee 是什么
- 参考 Gee 教程，在本仓库内手撸的最小 Web 框架，保持与 gin/现有代码完全隔离。  
- 目标：用最少依赖展示路由、中间件、模板、错误恢复、安全与限流的实现方式，便于学习框架原理。

## 运行最小 Demo
```bash
go run ./cmd/tinygee-demo --port :9999 --prom=true
# 访问：
# curl http://localhost:9999/
# curl http://localhost:9999/ping
# curl http://localhost:9999/metrics  # 若开启 --prom (expvar)
```

## 示例导航
- Day1 基础路由：`go run ./examples/tinygee/day1`
- Day3 动态路由：`go run ./examples/tinygee/day3`
- Day5 分组与中间件：`go run ./examples/tinygee/day5`
- Day6 模板 & 静态：`go run ./examples/tinygee/day6`
- Day7 Recover：`go run ./examples/tinygee/day7`
- Day9 JWT/RBAC：`go run ./examples/tinygee/day9auth`
- Day9 限流：`go run ./examples/tinygee/day9rl`
- Day10 Fullstack：`go run ./examples/tinygee/fullstack`（路由+模板+JWT/RBAC+限流+metrics）

## 测试
```bash
make tinygee-test        # go test ./tinygee/...
make tinygee-bench       # 预留基准，需补充 bench
```

## 关键特性
- 动态路由：支持 `:param` / `*filepath`
- 中间件链：Logger、Recover，可扩展
- 路由分组：前缀叠加 + 分组中间件
- 模板/静态：FuncMap + 模板渲染，静态文件服务
- 安全：JWT 验证、简单 RBAC 前缀控制
- 稳定性：Recover 防 panic；可选 Prometheus `/metrics`

## 后续待办
- Recover 配置化响应（已支持 JSON/文本，后续可拓展 HTML）
- `/metrics` httptest 覆盖
- Benchmarks：路由匹配、Middleware 开销
- Fullstack 示例：整合 JWT + 限流 + 模板 + Prom
