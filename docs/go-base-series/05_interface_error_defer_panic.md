# 05 interface、error、defer、panic

## 🎯 这一章解决什么问题

这章会真正把你从 Java 思维拉进 Go 风格。

需要解决：

1. Go 的接口为什么不需要 `implements`
2. Go 的错误处理为什么是 `if err != nil`
3. `defer` 到底什么时候执行
4. `panic/recover` 和 Java 异常到底有什么本质差异

## 🧩 最小代码

```go
package main

import (
"fmt"
"os"
)

type Notifier interface {
Notify(msg string) error
}

type Console struct{}

func (Console) Notify(msg string) error {
fmt.Println(msg)
return nil
}

func main() {
var n Notifier = Console{}
if err := n.Notify("hello"); err != nil {
    fmt.Println("notify failed:", err)
}

file, err := os.Open("README.md")
if err != nil {
    fmt.Println("open failed:", err)
    return
}
defer file.Close()
}
```

## 1️⃣ interface：隐式实现

```go
type Notifier interface {
Notify(msg string) error
}
```

只要某个类型实现了这个方法，它就自动满足接口。

```go
type Console struct{}

func (Console) Notify(msg string) error {
fmt.Println(msg)
return nil
}
```

这里没有写：

```go
implements Notifier
```

但它已经满足了接口。

这就是 Go 接口的核心：**隐式实现**。

## 2️⃣ error：错误是值，不是异常流

Go 没有 Java 那套 try-catch-finally 作为日常主流程。

Go 的主流写法是：

```go
file, err := os.Open("README.md")
if err != nil {
return
}
```

错误是普通返回值。

优点是：

- 调用点显式
- 错误路径清楚
- 不依赖隐式异常传播

## 3️⃣ `defer`：延迟到函数返回前执行

```go
defer file.Close()
```

含义：

- 现在先登记
- 等当前函数结束前再执行

最常见用途：

- 关闭文件
- 解锁 mutex
- 打印收尾日志

### 你可以先记一句

`defer` 很像“把收尾动作挂在函数出口”。

## 4️⃣ `panic/recover`

`panic` 表示程序进入不可恢复错误状态。

`recover` 只能在 `defer` 里拦截 panic。

但在日常业务逻辑里，Go 不推荐把它当成 try-catch 的平替。

正常业务错误，仍然应该用：

```go
return err
```

## 🆚 Java 开发怎么理解

### 接口

Java：

- 显式 `implements`

Go：

- 方法签名匹配就自动满足

### 错误

Java：

- 依赖异常传播机制

Go：

- 错误是返回值
- 调用方当场处理

### defer

可以类比 Java 的 `finally`，但不是完全等价。

`defer` 更像“提前登记一个退出时动作”。

## ⚠️ 注意点 / 易错点

### 1. 不要把 `panic` 当业务错误处理工具

`panic` 适合：

- 真正不该继续运行的状态
- 程序员错误
- 基础设施初始化失败且无法恢复

不适合：

- 用户输入错了
- 文件没找到
- 网络请求失败

### 2. `error` 可以是 `nil`

`error` 本身是接口类型。

判断错误要写：

```go
if err != nil {}
```

### 3. `defer` 是后进先出

```go
defer fmt.Println(1)
defer fmt.Println(2)
```

输出顺序是：

```text
2
1
```

### 4. 接口最好小而精

Go 风格通常不鼓励“大而全接口”。

小接口更灵活，也更适合测试替换。

### 5. 接口通常在“使用方”定义更合理

这是 Go 里非常典型的设计习惯。

## ▶️ 本章建议运行命令

```bash
go run ./examples/interfaces
go run ./examples/basics/error_wrap
go run ./examples/basics/defer_rules
go run ./examples/basics/interface_assert
```

## 📌 本章小结

你要记住：

1. Go 接口是隐式实现
2. `error` 是返回值，不是异常机制
3. `defer` 在函数退出前执行
4. `panic/recover` 不是日常业务错误流
5. Go 更偏好小接口和显式错误处理

## ⏭️ 下一章

下一章进入工程组织层：

- package
- module
- `go.mod`
- 目录结构
- `go run` / `go build` 的编译边界
