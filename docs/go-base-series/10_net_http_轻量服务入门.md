# 10 net/http 轻量服务入门

## 🎯 这一章解决什么问题

这章是整个 Go 基础阶段的收口。

目标是：

- 用标准库写一个最小 HTTP 服务
- 理解 `handler`、`ResponseWriter`、`Request`
- 返回 JSON
- 知道为什么学完这一层再去学 Gin 更稳

## 🧩 最小 HTTP 服务

```go
package main

import (
"encoding/json"
"net/http"
)

type Resp struct {
Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
_ = json.NewEncoder(w).Encode(Resp{Message: "hello go"})
}

func main() {
http.HandleFunc("/hello", helloHandler)
http.ListenAndServe(":8080", nil)
}
```

## 🔍 逐行解释

### `http.HandleFunc("/hello", helloHandler)`

表示注册一个路由：

- 路径：`/hello`
- 处理函数：`helloHandler`

### `func helloHandler(w http.ResponseWriter, r *http.Request)`

这个签名非常重要。

- `w`：往客户端写响应
- `r`：拿到请求信息

### `json.NewEncoder(w).Encode(...)`

直接把结构体编码成 JSON 写回客户端。

### `http.ListenAndServe(":8080", nil)`

启动一个 HTTP 服务，监听 8080 端口。

## 1️⃣ 为什么这已经足够作为“框架前基础”

因为 Gin、Echo 这些框架做的很多事，本质上都是在 `net/http` 之上做封装：

- 更方便的路由
- 更顺手的 JSON 绑定
- 更丰富的中间件机制
- 更好的参数提取体验

但只要你先理解了标准库这层，就不会把框架 API 当成黑盒魔法。

## 2️⃣ 一个稍微真实一点的例子

```go
package main

import (
"encoding/json"
"net/http"
)

type Todo struct {
ID    int    `json:"id"`
Title string `json:"title"`
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
todos := []Todo{
    {ID: 1, Title: "学习 Go"},
    {ID: 2, Title: "写一个 net/http 服务"},
}

w.Header().Set("Content-Type", "application/json")
_ = json.NewEncoder(w).Encode(todos)
}

func main() {
http.HandleFunc("/todos", todosHandler)
http.ListenAndServe(":8080", nil)
}
```

## 🆚 Java 开发怎么理解

可以先这样类比：

- `http.HandleFunc` 类似最基础的路由注册
- `http.Request` 类似请求对象
- `ResponseWriter` 类似响应输出对象
- 标准库 `net/http` 更接近 Servlet / 最小 Web 层，而不是 Spring MVC 的高层抽象

## ⚠️ 注意点 / 易错点

### 1. `ListenAndServe` 会阻塞当前 goroutine

所以通常它放在 `main()` 的最后。

### 2. 返回 JSON 时记得设置 Content-Type

```go
w.Header().Set("Content-Type", "application/json")
```

### 3. 写响应状态码要在写 body 前做

```go
w.WriteHeader(http.StatusBadRequest)
```

然后再写响应体。

### 4. 生产代码不要把所有错误都用 `_ =` 吞掉

教学 demo 可以简化，真实代码要处理编码、写响应等错误。

### 5. 标准库阶段先别急着做复杂路由

这一章的重点不是“造框架”，而是看懂 HTTP 服务最小闭环。

## 🔗 与仓库真实代码的连接点

当你学到这里，就可以去看：

- `cmd/todoapi`
- `internal/todo`
- `examples/java_compare/http_middleware`

重点看：

- handler 签名
- 路由注册
- JSON 编解码
- 中间件风格
- `context` 在请求链路中的传递

## ▶️ 本章建议运行命令

```bash
go run ./examples/java_compare/http_middleware
go run ./cmd/todoapi
```

然后测试：

```bash
curl http://localhost:8080/healthz
```

如果你已经熟悉登录流程，也可以继续看：

```bash
curl -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

## 📌 本章小结

你要记住：

1. `net/http` 是 Go Web 基础层
2. handler 负责处理请求并写回响应
3. `ResponseWriter` 和 `Request` 是最核心的两个对象
4. JSON 返回依赖 `encoding/json`
5. 学懂这一层，再看 Gin 会轻松很多

## ✅ 到这里你已经完成了什么

学到这里，Go 基础阶段可以认为已经完整闭环。

你已经过完：

- 语言基础
- 数据结构
- 方法与接口
- 错误处理
- 并发
- 常用标准库
- 工具链
- 轻量 HTTP 服务

下一阶段如果继续学，才适合进入：

- Gin / Echo / Fiber
- 配置管理
- 数据库访问
- 更完整的项目分层
