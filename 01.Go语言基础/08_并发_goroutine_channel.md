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

## 练习
- 使用 `context.WithTimeout` 包装 HTTP 请求
- 写一个简易 worker pool，限制并发数 5，处理字符串长度统计
