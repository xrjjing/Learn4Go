# Java 转 Go 快速学习指南

> 面向已经有 Java 开发经验、希望尽快系统进入 Go 的开发者。

## 🎯 核心结论

最快的学习方式，不是从头刷一遍纯语法书，而是按下面这条路线走：

1. 先建立 **Java → Go 的心智映射**
2. 用小程序快速过一遍核心语法
3. 尽快进入标准库、并发和工程化
4. 尽快写一个真实服务，而不是只停留在 demo
5. 全程坚持 **运行 → 修改 → 再运行**


## 🚀 推荐学习路线

### 第一阶段：先建立 Go 思维

先不求全，只求把主干打通：

- `package / import / func`
- `struct / method / interface`
- `slice / map / pointer`
- `error`
- `goroutine / channel / context`


### 第二阶段：用 Java 对照去理解 Go

建议直接建立下面这些映射：

- 类 → `struct`
- 方法 → 接收者方法
- `implements` → 隐式实现
- try-catch → `if err != nil`
- ThreadPool / Future → `goroutine + channel + WaitGroup`
- Spring 注入 → 手动组装依赖


### 第三阶段：直接写 HTTP 服务

优先掌握：

- `net/http`
- `encoding/json`
- `context`
- 路由、中间件、超时控制


### 第四阶段：再补工程化

重点掌握：

- `go mod`
- `go run`
- `go build`
- `go test`
- `go fmt`
- `go vet`


## 📚 结合本仓库的推荐学习顺序

先看这些文档：

1. `docs/Java对照Go速查表.md`
2. `docs/Go学习规划_Java开发者版.md`
3. `docs/Go学习资源.md`
4. `01.Go语言基础/`
5. `docs/Java对照练习.md`


## 🧭 具体学习顺序

### 1. 先学语言本身，不要先学框架

优先过完：

- 快速开始与基本语法
- 变量、常量、类型
- 流程控制
- 函数与错误处理
- 数组、切片、Map
- 结构体、方法、组合
- 接口与多态
- 并发：goroutine 与 channel
- 包与模块管理
- 测试与基准

对应目录：

```text
01.Go语言基础/
```


### 2. 用 Java 对照示例去建立直觉

建议优先跑这些：

```bash
go run ./examples/java_compare/interface_poly
go run ./examples/java_compare/concurrency
go run ./examples/java_compare/context_timeout
go run ./examples/java_compare/ticker_rate_limit
go run ./examples/java_compare/http_middleware
go run ./examples/java_compare/syncmap
```

这些示例适合建立下面的认知：

- 接口如何实现
- Go 并发怎么写
- 超时和取消怎么做
- 中间件怎么组织
- 并发安全如何处理


### 3. 尽快掌握标准库

Go 的核心竞争力不只是语法，而是标准库。

建议优先掌握：

- `fmt`
- `time`
- `strings`
- `os`
- `io`
- `encoding/json`
- `net/http`
- `context`
- `sync`


### 4. 尽快进入真实项目

语法过一遍后，不要停留在玩具示例。

建议直接运行：

```bash
go run ./cmd/todoapi
```

然后重点阅读：

- `cmd/todoapi`
- `internal/todo`

重点关注：

- 程序入口如何组织
- handler 怎么写
- service/store 如何拆分
- 错误如何返回
- JSON 如何处理
- 依赖如何组装


## 💡 Java 转 Go 最重要的 6 个思维转换

### 1. 不要总想着 class

Go 没有 Java 那种“一切都要放进类里”的习惯。

更常见的做法是：

- `struct` 放数据
- 方法挂在接收者上
- 普通函数直接存在于 package 下


### 2. 不要期待继承

Go 没有 `extends`。

核心思路是：

- 组合
- 嵌入
- 小接口


### 3. 不要依赖异常流控

Go 的错误处理是显式的：

```go
v, err := xxx()
if err != nil {
    return err
}
```

刚开始可能会觉得啰嗦，但这是 Go 的核心风格。


### 4. 不要把接口理解成“实现方声明”

Java：

```java
class A implements B
```

Go 不是这样。

Go 是：**方法签名匹配，就自动满足接口**。


### 5. 不要一开始就依赖重量框架

Java 很容易先学 Spring Boot。

Go 更推荐：

1. 先学语言
2. 再学标准库
3. 最后再学框架


### 6. 并发不是加分项，而是主能力

必须尽早掌握：

- goroutine
- channel
- `select`
- `context`
- `sync.WaitGroup`
- `sync.Mutex`


## 📅 一个适合 Java 开发者的 14 天启动计划

### 第 1-3 天

- 阅读 `docs/Java对照Go速查表.md`
- 学 `01.Go语言基础` 前 4 章
- 自己手写 5 个小 demo：
  - hello world
  - 变量与 `:=`
  - `if / for`
  - 函数多返回值
  - `error`


### 第 4-6 天

- 学切片、map、struct、method、interface
- 跑这些示例：

```bash
go run ./examples/java_compare/interface_poly
go run ./examples/java_compare/bufio_wordcount
```


### 第 7-10 天

- 学 goroutine、channel、WaitGroup、context
- 跑这些示例：

```bash
go run ./examples/java_compare/concurrency
go run ./examples/java_compare/context_timeout
go run ./examples/java_compare/ticker_rate_limit
```


### 第 11-14 天

- 学 `net/http`、JSON、模块、测试
- 跑这些命令：

```bash
go run ./cmd/todoapi
GOCACHE=$(pwd)/.gocache go test ./...
```

完成这一步后，基本就已经不再是“刚接触 Go”的状态了。


## 🛠️ 每天建议的学习节奏

每天 1.5 ~ 2 小时，建议这样安排：

1. 15 分钟：看一个小知识点
2. 30 分钟：运行仓库里的示例
3. 30 分钟：自己改一版
4. 15 分钟：总结它和 Java 的最大差异

这个节奏通常比单纯看视频更快、更稳。


## 📌 优先级最高的 8 个知识点

先把下面这些打透：

- `package`
- `func`
- `struct`
- `interface`
- `slice`
- `error`
- `goroutine`
- `context`

这 8 个打通，Go 的主体就通了。


## ⚠️ 刚开始不建议过早深挖的内容

前期先别把太多时间花在：

- 反射
- `unsafe`
- 复杂泛型
- 过多三方框架
- 过早看微服务全家桶
- 复杂设计模式迁移

先把 Go 的简单、显式、直接的风格吃透。


## ✨ 最重要的一句话

**不要把 Go 学成“没有 Spring 的 Java”。**

要主动接受 Go 的风格：

- 少魔法
- 少抽象
- 显式
- 简单
- 标准库优先


## ✅ 建议的下一步

如果现在就开始，推荐直接按这个顺序动手：

```bash
# 1. 先看 Java 对照
sed -n '1,200p' docs/Java对照Go速查表.md

# 2. 再跑几个 Java 对照示例
go run ./examples/java_compare/interface_poly
go run ./examples/java_compare/concurrency
go run ./examples/java_compare/context_timeout

# 3. 最后进入真实服务
go run ./cmd/todoapi
```

学 Go 的关键不是一次性把所有知识点都看完，而是尽快形成：

**能看懂 → 能运行 → 能修改 → 能自己写**
