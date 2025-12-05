# Go 语言基础资源与练习建议

## 官方与中文教程

- Go 官方文档入口：语言规范、标准库、工具链、内存模型，权威且持续更新，包含 Tour、Effective Go、How to write Go code 等。citeturn1search0turn1search4
- 菜鸟教程 Go 专栏：涵盖安装、基础语法、指针、泛型、循环等章节，示例简洁，便于快速对照练习。citeturn0search0turn0search2turn0search3turn0search7

## 配套内部文档与示例

- 工具链：`docs/Go工具链.md`，示例 `examples/basics/go_generate/`
- IO：`docs/IO操作指南.md`，示例 `examples/io/*`
- 基础速查：`docs/Go基础速查.md`
- 测试/pprof：`docs/性能分析pprof.md`，示例 `examples/testing/benchmark_pprof`

## 建议的学习顺序（面向 Java 背景）

1. 语法快速扫盲：变量/常量、零值、短变量声明、作用域与可见性。
2. 数据结构：数组 vs 切片（容量/扩容）、map 的 nil 约束。
3. 函数与错误处理：多返回值、命名返回值、`if err != nil`，对比 try-catch。
4. 方法与接口：指针接收者 vs 值接收者、隐式接口实现。
5. 并发基础：goroutine、channel、`select`，对比 Java 线程池。
6. 标准库速览：`net/http`（内置服务器）、`context`（超时/取消）、`encoding/json`、`sync` 系列。
7. 工具链：`go test`、`go vet`、`go fmt`、`go build`、`go mod`。

## 配套练习（仓库内示例）

- 本地运行：`examples/java_compare/*`（并发/超时/中间件/pprof/sync.Map 等对照练习）。
- 快速体验网络特性：`examples/java_compare/httptrace`（需外网 httpbin），`examples/java_compare/http_middleware`（本地）。
- Ticker 限速与 context 取消：`examples/java_compare/ticker_rate_limit`、`context_timeout`。
- 词频统计：`examples/java_compare/bufio_wordcount`，体会 bufio 与字符串处理。

## 进阶阅读

- 官方内存模型、go tool 文档、包文档（建议按需查阅标准库包说明）。citeturn0search2

## 小贴士

- 以运行示例为驱动：每个示例先 `go run` 观察输出，再修改参数体验行为变化。
- 将 Java 对应实现放在旁边对照，关注错误处理、接口实现和并发模型的差异。
- 结合 `pprof_server` 示例，使用 `go tool pprof` 打开本地火焰图，理解性能调优入口。\*\*\*
