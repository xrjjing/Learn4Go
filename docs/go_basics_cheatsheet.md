# Go 语言基础速查（含练习指引）

## 指针与值语义
- 规则：小对象/不可变可用值接收者；需要共享/修改状态用指针接收者。
- 练习：`examples/basics/pointer_receiver`，对比值/指针接收者的行为和逃逸。

## 零值语义
- 大多数类型零值可用；`map`/`chan`/`slice` 需 make。
- 练习：`examples/basics/map_basics`（nil map 写入 panic 演示与修复）。

## 切片与扩容
- len/cap、共享底层数组、切片追加导致原数据被改写。
- 练习：`examples/basics/slice_aliasing`，观察两个切片共享底层的效果。

## defer 细节
- LIFO、参数按值捕获、闭包变量按引用捕获；循环内 defer 的开销。
- 练习：`examples/basics/defer_rules`。

## 错误处理
- 哨兵错误、wrapping：`fmt.Errorf("…: %w", err)`；`errors.Is/As`。
- 练习：`examples/basics/error_wrap`。

## goroutine 与 channel
- 泄露防治：使用 ctx 或关闭 channel；缓冲/无缓冲差异；select 超时。
- 练习：`examples/basics/channel_timeout`（select + time.After）。

## context 传递
- 约定：第一个参数 ctx；避免存业务大对象；`WithTimeout/WithCancel`。
- 练习：`examples/basics/context_cancel`。

## 工具链最小集
- `go fmt`、`go vet`、`go test`、`go mod tidy`；逃逸分析：`go test -gcflags=-m ./...`（可选）。

## 如何使用
1. 按目录运行：`go test ./examples/basics/...` 可直接看到断言。  
2. 修改练习代码中的变量/容量/超时参数，观察行为变化。  
3. 对照 `docs/go_basics_resources.md` 的外部教程获取更详细解释。***
