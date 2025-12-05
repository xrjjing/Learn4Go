# Go 语言基础差距梳理（对标 c.biancheng.net/golang/）

本仓库已有：

- 语法入门与对照：`01.Go语言基础` 文档、`examples/basics/`（指针/切片/map/defer/错误/并发/ctx）
- 并发与网络：channel/select、workerpool、gateway/gRPC 示例
- 工程化：go mod、脚本、TinyGee 手撸框架、JWT/RBAC/限流、模板/静态、Redis/MinIO、Docker 部署

参考 biancheng.net 目录，尚需补充/加强的基础点：

## 1. 环境与工具链

- GOPATH、GOROOT、GOBIN 详解及 `go env` 输出解读（含 Go Modules 时代区别）。
- go get / go mod download / go generate 用法与示例。
  > 计划：新增文档 `docs/Go工具链.md`，配套示例 `examples/basics/go_generate/`。

## 2. 语言基础细节

- 常量与 `iota` 进阶用法（位掩码、自增技巧）。
- 数值类型边界/溢出、类型转换规则。
- 字符串与 rune/byte 细节（UTF-8、切片风险）。
- 运算符与优先级、短路逻辑。
  > 计划：在 `examples/basics/` 增加 `iota_bits`、`rune_utf8`、`numeric_limits`。

## 3. 控制流与函数

- switch 多分支、fallthrough、类型 switch；for-range 拆解。
- 可变参数/闭包/递归/延迟调用性能提示。
  > 计划：补 `examples/basics/switch_fallthrough`、`closure_defer_cost`。

## 4. 结构体与接口

- 组合/匿名字段、方法集规则；接口断言/类型转换 panic 安全写法。
  > 计划：补文档片段 `docs/Go基础速查.md` 对应条目，并增加 `examples/basics/interface_assert`.

## 5. 错误与异常

- 自定义错误、哨兵错误 vs wrapped、panic/recover 最佳实践（已覆盖基础，可补案例：多层 wrap 与 unwrap）。

## 6. 文件与 IO（biancheng 重点）

- bufio/文件读写、复制、锁；zip/tar/gzip；JSON/XML 读写。
  > 计划：新增 `examples/io/` 系列：`bufio_rw`、`copy_file`、`zip_basic`、`json_xml`. 文档 `docs/IO操作指南.md`。

## 7. 并发与定时

- time.Ticker/Timer、select 超时（部分已有）；等待组、互斥/读写锁/原子操作讲解。
  > 计划：`examples/concurrency/waitgroup_mutex_atomic`.

## 8. 网络与标准库

- TCP/UDP socket 基础；net/http client/server 更详尽示例；context 取消 HTTP。
  > 计划：`examples/network/tcp_echo`、`network/udp_echo`、`network/http_client_ctx`.

## 9. 测试与性能

- go test 表格驱动（已有）、benchmark + pprof 使用。
  > 计划：`examples/testing/benchmark_pprof`，文档 `docs/性能分析pprof.md`.

## 10. 项目结构与包管理

- GOPATH 模式 vs Modules 模式下的目录差异；常见分层（MVC 等）对比。
  > 计划：补充 `docs/架构设计.md` 小节，引用 TinyGee/Gateway 作为示例。

## 11. 泛型（Go1.18+）

- 虽 biancheng 旧文未覆盖，但本项目可增加最小泛型示例。
  > 计划：`examples/basics/generic_stack`.

## 执行优先级（建议）

1. IO + 工具链（go get/go mod/go generate）✅
2. 数值/字符串/iota/rune 细节 ✅
3. 网络 TCP/UDP/HTTP 客户端 ✅
4. 测试/benchmark/pprof ✅
5. 泛型示例 ✅

完成后同步：

- README“基础学习资源”列表与 `docs/前端使用指南.md` 示例导航。
- 门户示例过滤可追加“IO/网络/泛型”分组标签（可选）。

引用：c.biancheng.net/golang 目录及相关章节。citeturn0search0
