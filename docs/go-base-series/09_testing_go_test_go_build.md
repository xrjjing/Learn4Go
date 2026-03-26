# 09 测试、构建、工具链

## 🎯 这一章解决什么问题

学语言只会写 demo 不够，还要知道怎么验证、构建和交付。

这章解决：

1. `go run`、`go build`、`go test` 分别干什么
2. `go build -o` 应该怎么用
3. 为什么按文件编译和按目录编译差别很大
4. `go fmt`、`go vet` 该怎么看待

## 1️⃣ `go run`

```bash
go run hello.go
go run ./examples/hello
```

作用：

- 临时编译
- 直接运行
- 不保留最终二进制

很适合学习阶段和快速验证。

## 2️⃣ `go build`

```bash
go build
go build -o hello hello.go
go build -o ./bin/app ./cmd/todoapi
```

作用：

- 编译成可执行文件或产物
- 不自动运行

### `-o` 是什么

`-o` 是指定输出文件名或路径。

例如：

```bash
go build -o hello hello.go
```

表示把 `hello.go` 编译成 `hello` 这个可执行文件。

## 3️⃣ `go test`

Go 的测试是工具链内建能力。

- 测试文件后缀：`_test.go`
- 测试函数签名：`func TestXxx(t *testing.T)`

常见命令：

```bash
go test ./...
go test -bench=. ./examples/testing
```

## 4️⃣ `go fmt`

Go 鼓励把代码风格交给工具，不做太多个性化争论。

```bash
go fmt ./...
```

## 5️⃣ `go vet`

`go vet` 不是格式化工具，而是静态检查工具。

它会帮你发现一些可疑代码，比如：

- 格式化参数和类型不匹配
- 不太合理的代码模式

## 🧩 一个最小测试示例

```go
package calc

func Add(a, b int) int {
return a + b
}
```

```go
package calc

import "testing"

func TestAdd(t *testing.T) {
if got := Add(2, 3); got != 5 {
    t.Fatalf("got %d, want 5", got)
}
}
```

## 🆚 Java 开发怎么理解

- `go test` 类似 JUnit 执行入口，但更内建
- `go build` 类似 Maven/Gradle 的编译产物阶段，但更轻
- `go fmt` 类似统一格式化插件，但在 Go 世界里几乎是默认规则
- `go vet` 更像轻量静态检查

## ⚠️ 注意点 / 易错点

### 1. `go build -o` 后面必须跟输出目标

不能只写：

```bash
go build -o
```

### 2. 按文件构建和按目录构建差别很大

```bash
go build -o hello hello.go  # 只编一个文件
go build -o hello .         # 编当前整个包
```

这点正是你之前遇到 `main redeclared` 的根源之一。

### 3. 学习时优先 `go run`，交付时再 `go build`

这样心智更清晰。

### 4. 测试默认按 package 维度跑

所以目录组织直接影响测试体验。

### 5. 格式化是工具责任，不要手工死磕对齐

写完就 `go fmt`，让工具统一。

## ▶️ 本章建议运行命令

```bash
go run ./examples/testing
GOCACHE=$(pwd)/.gocache go test ./...
go test -bench=. ./examples/testing
make fmt vet test
```

## 📌 本章小结

你要记住：

1. `go run` 适合快速验证
2. `go build` 负责产出可执行文件
3. `-o` 用来指定输出路径
4. `go test` 是内建测试能力
5. `go fmt` 和 `go vet` 是基础工程习惯

## ⏭️ 下一章

最后一章开始把前面知识收口到一个真实方向：

- 用 `net/http` 写一个轻量服务
- 理解 handler、路由、JSON 响应、状态码
