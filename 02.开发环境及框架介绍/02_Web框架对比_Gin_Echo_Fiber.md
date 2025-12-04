# Web 框架对比：Gin / Echo / Fiber（概览）

| 维度 | Gin | Echo | Fiber |
| --- | --- | --- | --- |
| 性能 | 高 | 高 | 极高（基于 fasthttp） |
| 中间件 | 丰富、社区成熟 | 简洁、路由清晰 | 语法类似 Express，学习曲线低 |
| 特点 | 上手快、文档多 | API 设计简洁 | 速度优先，但与 net/http 不完全兼容 |

## 选择建议
- 入门与社区资料：Gin
- 极简与更少依赖：Echo
- 追求性能或 Node 习惯：Fiber

## 标准库起步
在无法拉取依赖时，可先用 `net/http` + 自写中间件模式：
```go
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

## 练习
- 用 `net/http` 编写一个 `/healthz` 与 `/echo` 接口
- 对比同样逻辑在 Gin/Echo/Fiber 的写法，记录代码行数与可读性
