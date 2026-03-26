# 07 并发：goroutine、channel、context

## 🎯 这一章解决什么问题

这章是 Go 的核心竞争力之一。

需要解决：

1. goroutine 和 Java Thread / 线程池有什么区别
2. channel 到底是队列、管道，还是同步器
3. `select` 有什么用
4. `context` 为什么在 Go 里这么重要

## 🧩 最小代码

```go
package main

import (
"fmt"
"time"
)

func worker(ch chan string) {
time.Sleep(time.Second)
ch <- "done"
}

func main() {
ch := make(chan string)
go worker(ch)
msg := <-ch
fmt.Println(msg)
}
```

## 1️⃣ goroutine：非常轻量的并发执行单元

```go
go worker(ch)
```

`go` 关键字会启动一个新的 goroutine。

你可以先把 goroutine 理解成：

- 比线程轻得多
- 创建成本低
- Go 运行时负责调度

这不是“你手动 new 一个 Thread”的思路。

## 2️⃣ channel：通过通信共享数据

```go
ch := make(chan string)
ch <- "done"
msg := <-ch
```

channel 是 Go 并发里非常标志性的原语。

你可以先把它理解成：

- 带同步语义的通道
- 可以传值
- 也可以当作 goroutine 间协调工具

### 无缓冲 channel

```go
ch := make(chan string)
```

发送和接收要配对，否则会阻塞。

### 有缓冲 channel

```go
ch := make(chan string, 3)
```

有一定暂存能力，但也不是无限队列。

## 3️⃣ `select`：同时监听多个 channel

```go
select {
case msg := <-ch:
fmt.Println(msg)
case <-time.After(time.Second):
fmt.Println("timeout")
}
```

这个能力非常重要，用来：

- 等待多个事件
- 做超时控制
- 配合 context 取消任务

## 4️⃣ `context`：控制超时、取消、请求边界

在 Go 里，尤其是 HTTP / RPC / DB 调用链，`context` 非常重要。

你可以把它理解成：

- 当前请求或任务的上下文
- 可以携带取消信号
- 可以设置截止时间或超时

常见写法：

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
```

## 🆚 Java 开发怎么理解

### goroutine

不要把它简单等同成 `new Thread()`。

更接近：

- 轻量任务
- 由运行时调度
- 配合 channel 和 context 使用

### channel

可以把它类比成：

- `BlockingQueue`
- `Future` 的结果通道
- 带同步语义的消息通道

但 channel 又比这些更基础。

### context

可以类比 Java 中“任务取消 + 超时令牌 + 请求上下文”的组合，但 Go 有统一约定，很多标准库和框架都会接收 `context.Context`。

## ⚠️ 注意点 / 易错点

### 1. goroutine 不是越多越好

它很轻量，但不是没成本。

### 2. 没人接收的 channel 发送会阻塞

```go
ch := make(chan int)
ch <- 1 // 如果没有接收者，会卡住
```

### 3. 从没人发送的 channel 接收也会阻塞

```go
x := <-ch
```

如果没人发送，会一直等。

### 4. 不正确关闭 channel 会 panic

一般原则：

- 由发送方关闭
- 不要重复关闭

### 5. `context` 创建后记得 cancel

特别是 `WithTimeout` / `WithCancel` 创建出来的上下文，通常要：

```go
defer cancel()
```

### 6. 不要把 context 当万能参数包

它主要是：

- 超时
- 取消
- 请求范围元信息

不是拿来装业务字段的“大口袋”。

## ▶️ 本章建议运行命令

```bash
go run ./examples/concurrency
go run ./examples/workerpool
go run ./examples/java_compare/concurrency
go run ./examples/java_compare/context_timeout
go run ./examples/java_compare/ticker_rate_limit
go run ./examples/basics/channel_timeout
go run ./examples/basics/context_cancel
```

## 📌 本章小结

你要记住：

1. goroutine 是轻量并发单元
2. channel 用来通信和同步
3. `select` 用来同时等待多个事件
4. `context` 用来传递超时和取消边界
5. Go 并发强调“通过通信共享内存”

## ⏭️ 下一章

下一章进入最常用标准库：

- `time`
- `strings`
- `os`
- `io`
- `encoding/json`
