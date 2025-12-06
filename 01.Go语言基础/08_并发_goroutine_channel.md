# 08 并发：goroutine 与 channel

## 关键点

- goroutine：`go f()` 创建轻量线程
- channel：无缓冲/带缓冲，同步与通信
- select：多路复用；结合 `context` 控制取消/超时

## 示例

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 5)
    results := make(chan int, 5)
    for w := 1; w <= 2; w++ { go worker(w, jobs, results) }
    for j := 1; j <= 5; j++ { jobs <- j }
    close(jobs)
    for a := 0; a < 5; a++ { fmt.Println(<-results) }
}
```

## 高级技巧：原子操作

Go 提供了 `sync/atomic` 包实现无锁的原子操作，性能优于互斥锁：

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

// 计数器并发安全实现
type Counter struct {
    value atomic.Int64
}

func (c *Counter) Increment() int64 {
    return c.value.Add(1)
}

func (c *Counter) Get() int64 {
    return c.value.Load()
}

func main() {
    counter := &Counter{}
    var wg sync.WaitGroup

    // 1000 个 goroutine 并发递增
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()
    fmt.Printf("Final count: %d\n", counter.Get()) // 输出：Final count: 1000
}
```

### 实战案例：项目中的 ID 生成器

在 `internal/todo/store.go` 中，我们使用原子操作实现并发安全的 ID 生成器：

```go
type Store struct {
    mu     sync.Mutex
    items  map[int]Todo
    nextID atomic.Int64  // 原子递增 ID
}

func (s *Store) Create(title string, userID uint) (Todo, error) {
    id := int(s.nextID.Add(1))  // 原子递增，无需加锁

    t := Todo{
        ID:        id,
        UserID:    userID,
        Title:     title,
        CreatedAt: time.Now(),
    }

    s.mu.Lock()
    s.items[id] = t  // 只在修改 map 时加锁
    s.mu.Unlock()

    return t, nil
}
```

**优点**：

- ✅ 并发安全（原子操作保证）
- ✅ 高性能（无锁操作，0 次堆分配）
- ✅ ID 唯一（严格递增）

**性能对比**：

- 使用 `sync.Mutex` 锁整个 Create：~500 ns/op
- 使用 `atomic.Int64`：~284 ns/op（**快 43%**）

## 练习

- 使用 `context.WithTimeout` 包装 HTTP 请求
- 写一个简易 worker pool，限制并发数 5，处理字符串长度统计
- 使用 `sync/atomic` 实现一个线程安全的统计器（支持递增、递减、获取值）
