# Java vs Go 速查对照

| 主题 | Java | Go |
| --- | --- | --- |
| 包管理 | Maven/Gradle | Go Modules (`go mod tidy`) |
| 可见性 | public/protected/private | 首字母大写导出，小写包内可见 |
| 异常 | checked/unchecked | panic/recover，错误为值 `error` |
| 并发 | Thread/Executor/Future | goroutine + channel + `context.Context` |
| 泛型 | Java 泛型 (擦除) | Go1.18+ type parameters |
| 集合 | List/Map/Set | 切片/数组/map，内置迭代 | 
| 测试 | JUnit/TestNG | `go test`，基准 `go test -bench` |
| Web | Spring Boot | 标准库 `net/http`，生态 Gin/Echo/Fiber |
| 序列化 | Jackson/Gson | `encoding/json`，结构体 tag |
| 依赖注入 | Spring | wire/google injector 等，或手写构造函数 |

## 关键差异
- **接口**：Go 接口为结构子类型，无显式 `implements`，解耦更强。
- **并发模型**：goroutine 更轻量；channel 负责同步与通信；`select` 等待多路事件。
- **错误处理**：返回 `error` 值，按需包裹（fmt.Errorf/"%w"）。
- **初始化顺序**：包级变量 → init() → main()；避免复杂 init 副作用。

## 常用对照代码
```go
// Java Optional vs Go 零值+指针
func find(id int) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid id")
    }
    return &User{ID: id, Name: "demo"}, nil
}
```

```go
// Java synchronized vs Go channel/Mutex
var mu sync.Mutex
mu.Lock()
// 临界区
mu.Unlock()
```

```go
// CompletableFuture vs goroutine + channel
ch := make(chan int)
go func() {
    ch <- compute()
}()
result := <-ch
```
