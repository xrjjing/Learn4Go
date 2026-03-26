# Go 基础连载总纲

> 面向 Java 开发者的 Go 基础连载讲义。
> 学到这里为止，你应该能独立写小型 CLI、读懂标准库代码、写出轻量 HTTP 服务，然后再进入 Gin / GORM / gRPC。

## 🎯 这套连载解决什么问题

你现在的问题不是“找不到 Go 资料”，而是：

- 资料太散，缺少主线
- Java 思维容易带偏对 Go 的理解
- 单看语法会，会抄，但不会判断
- 经常在命令、作用域、包结构这些地方踩坑

这套连载的目标是把 Go 基础阶段讲成一条连续主线：

1. 先理解程序怎么启动
2. 再理解变量、函数、作用域和错误处理
3. 再进入数据结构、结构体、接口、并发
4. 再补标准库、工具链、测试
5. 最后用 `net/http` 写一个轻量服务收口

## 🧭 学习边界

这套连载 **讲到框架前一层** 为止，包含：

- 语言基础：变量、函数、流程控制、数据结构、方法、接口
- 工程基础：包、模块、构建、测试、常见命令
- 并发基础：goroutine、channel、context
- 标准库基础：fmt、time、strings、os、io、json、http
- 最终收口：使用 `net/http` 写轻量 HTTP 服务

明确 **不进入**：

- Gin / Echo / Fiber
- GORM / Ent
- gRPC
- 微服务拆分、配置中心、服务治理
- 反射、unsafe、复杂泛型

## 📚 阅读顺序

```text
01 程序入口、变量、函数
02 流程控制、作用域、短变量声明
03 数组、切片、map
04 struct、方法、指针、组合
05 interface、error、defer、panic
06 包、模块、工程组织
07 并发：goroutine、channel、context
08 常用标准库：io/json/time/strings/os
09 测试、构建、工具链
10 net/http 轻量服务入门
```

## 🧱 每章固定结构

每一章默认都包含：

- 这一章要解决什么问题
- 最小可运行代码
- 逐行解释
- Java 开发怎么理解
- 注意点 / 易错点 / 命令边界
- 可直接运行的命令
- 下一章衔接

也就是说，后续你就算只说一句“继续”，我也知道应该往哪一章接着讲。

## ⚠️ 这套连载会主动写进去的“联想注意点”

以后凡是这些高频坑，我不会等你追问，而会直接写进章节：

- `:=` 是声明，不是普通赋值
- `=` 是给已有变量赋值
- `err` 在不同作用域里可能被遮蔽
- 同目录多个 `main()` 会触发 `main redeclared`
- `go run file.go` 和 `go run .` 的编译范围不同
- `go build -o app file.go` 和 `go build -o app .` 的含义不同
- `fmt.Printf(printX())` 这类“无返回值当成值”的错误
- `map` 零值可读不可写
- `slice` 共享底层数组带来的连带修改
- `panic/recover` 不是 Go 版 try-catch
- 中文目录名、非法包路径、IDE/gopls 误报的常见原因

## 🔗 与仓库现有资料的关系

这套连载不是替换仓库原有资料，而是把原有资料重新串起来。

可并行参考：

- `Java转Go快速学习指南.md`：总的学习路线
- `docs/Java对照Go速查表.md`：Java → Go 心智迁移
- `docs/Go学习资源.md`：现有文档与示例索引
- `01.Go语言基础/`：更精简的原始章节
- `docs/Java对照练习.md`：仓库内对照练习入口

## 🗺️ 每章与现有示例的映射

```text
01 -> examples/hello, examples/variables, examples/functions
02 -> examples/controlflow, demo/
03 -> examples/collections, examples/basics/map_basics, examples/basics/slice_aliasing
04 -> examples/structs, examples/basics/pointer_receiver
05 -> examples/interfaces, examples/basics/error_wrap, examples/basics/defer_rules
06 -> examples/packages, docs/Go工具链.md
07 -> examples/concurrency, examples/workerpool, examples/java_compare/concurrency,
      examples/java_compare/context_timeout, examples/basics/channel_timeout,
      examples/basics/context_cancel
08 -> examples/java_compare/bufio_wordcount, examples/java_compare/urlencode,
      examples/basics/rune_utf8
09 -> examples/testing, examples/testing/benchmark_pprof
10 -> cmd/todoapi, examples/java_compare/http_middleware
```

## ✅ 学完这套连载后你应该能做到什么

学完以后，你应该具备这些能力：

- 看懂绝大多数基础 Go 代码
- 独立写一个 CLI 小工具
- 理解 `if err != nil`、多返回值、接口和组合
- 理解 Go 最常见的并发写法
- 独立写一个基于 `net/http` 的小服务
- 看 Gin / GORM / gRPC 时不再是“硬记 API”

## 🚀 推荐使用方式

最推荐的学习节奏：

1. 读一章
2. 跑里面的命令
3. 自己把示例改一版
4. 记录“这和 Java 最不一样的点”
5. 再继续下一章

如果中途有新问题，比如：

- 为什么这里报错？
- 这里为什么不能这么写？
- 这个地方 Java 要怎么类比？

这些问题后续都可以继续沉淀回章节里。

## 📌 当前连载目录

- `docs/go-base-series/01_程序入口_变量_函数.md`
- `docs/go-base-series/02_流程控制_作用域_短变量声明.md`
- `docs/go-base-series/03_数组_切片_map.md`
- `docs/go-base-series/04_struct_方法_指针_组合.md`
- `docs/go-base-series/05_interface_error_defer_panic.md`
- `docs/go-base-series/06_包_模块_工程组织.md`
- `docs/go-base-series/07_并发_goroutine_channel_context.md`
- `docs/go-base-series/08_io_json_time_strings_os.md`
- `docs/go-base-series/09_testing_go_test_go_build.md`
- `docs/go-base-series/10_net_http_轻量服务入门.md`
