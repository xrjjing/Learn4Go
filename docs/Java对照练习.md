# Java 背景的 Go 练习清单

## 标准库与语法对照
- `[本地]` `examples/java_compare/urlencode`：`net/url` vs Java `URLEncoder`  
- `[本地]` `examples/java_compare/bufio_wordcount`：`bufio.Scanner` 词频统计 vs `BufferedReader`  
- `[本地]` `examples/java_compare/interface_poly`：隐式接口实现 vs `implements`

## 并发与超时
- `[本地]` `examples/java_compare/concurrency`：`WaitGroup+channel` 并发抓取 vs `ExecutorService+Future`  
- `[本地]` `examples/java_compare/context_timeout`：`context.WithTimeout` 取消任务 vs `Future.get(timeout)`  
- `[本地]` `examples/java_compare/ticker_rate_limit`：`time.Ticker` 限速 vs `ScheduledExecutorService`  
- `[🌐 httpbin]` `examples/java_compare/httptrace`：`httptrace` 观测 DNS/连接/首字节 vs Java HttpClient 监听器  
- `[本地]` `examples/java_compare/http_middleware`：`net/http` 中间件链 vs Servlet Filter  
- `[本地]` `examples/java_compare/pprof_server`：内置 pprof 采样 vs Flight Recorder/VisualVM  
- `[本地]` `examples/java_compare/syncmap`：`sync.Map` 读多写少 vs ConcurrentHashMap

## 建议练习步骤
1. 逐个运行示例，观察输出与 Java 类比。  
2. 修改参数（URL、超时、ticker 间隔）体会行为变化。  
3. 将 `httptrace` 与 `context.WithTimeout` 组合，感受超时对链路的影响。  
4. 在 `bufio_wordcount` 中加入停用词过滤，练习字符串处理。  
5. 将 `interface_poly` 扩展为 `Storage` 的 Redis/File/Memory 多实现，对比 Java 的依赖注入。

## 运行命令示例
```bash
go run ./examples/java_compare/urlencode
go run ./examples/java_compare/bufio_wordcount
go run ./examples/java_compare/concurrency
go run ./examples/java_compare/context_timeout
go run ./examples/java_compare/interface_poly
go run ./examples/java_compare/ticker_rate_limit
go run ./examples/java_compare/httptrace
go run ./examples/java_compare/http_middleware
go run ./examples/java_compare/pprof_server
go run ./examples/java_compare/syncmap
```
