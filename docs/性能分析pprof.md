# 测试与性能分析速览

## benchmark
- 写法：`func BenchmarkXxx(b *testing.B)`；使用 `b.N` 控制循环。
- 运行：`go test -bench=. ./examples/testing/benchmark_pprof`。

## pprof
- CPU：`go test -bench=. -benchmem -cpuprofile=cpu.out ./examples/testing/benchmark_pprof`  
- 内存：`-memprofile=mem.out`  
- 查看：`go tool pprof cpu.out`，命令 `top`、`web`、`list`。

## flamegraph（可选）
- `go tool pprof -http=:8081 cpu.out` 打开 Web UI。

## 小贴士
- 基准前关闭其他负载，使用 GOMAXPROCS/环境一致。  
- 避免在 benchmark 中打印或分配大对象；使用 `b.ReportAllocs()` 查看分配。
