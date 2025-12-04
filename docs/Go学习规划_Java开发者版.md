# Go 学习规划（面向 Java 开发者）

> 节奏参考原 Python learn 仓库：语言基础 → 工程化 → 框架 → 实战递进。每周至少完成 2-3 章基础 + 1 个小练习。

## 阶段 0：准备
- 安装 Go 1.20+，配置 `GOPATH`（若需要）与 `GOMODCACHE`
- 工具：gofmt / gofmt -w、go test、go tool pprof、delve 调试
- 阅读 `docs/Java_vs_Go_CheatSheet.md`，建立心智对照

## 阶段 1：语言基础（1-2 周）
- 章节对应：`01.Go语言基础/01~07`
- 重点：值类型/引用类型、切片与 map 语义、defer/panic/recover、interface + 嵌入、并发原语(goroutine/channel)
- 练习：
  - 写一个 Fibonacci 生成器（channel 版本）
  - 用 `context.Context` 控制超时

## 阶段 2：工程化与框架（1 周）
- 章节对应：`02.开发环境及框架介绍`
- 目标：模块管理、配置、日志、分层结构、Web 框架对比（Gin/Echo/Fiber）、gRPC 基础
- 练习：
  - 将 HTTP mini 服务拆为 handler/service/repo 三层
  - 使用 `go test -race` 检查简单并发示例

## 阶段 3：项目实战（2-3 周）
- 章节对应：`03.项目实战`
- 路线：命令行工具 → 日志处理 → HTTP API → 并发/worker → gRPC/gateway
- 练习：
  - 为 CLI 增加 `--dry-run` 与 `--apply` 模式
  - 为 mini TODO API 增加中间件：请求日志、统一响应
  - 实现一个内存任务队列，带超时与取消

## 阶段 4：巩固与扩展
- 编写基准测试与 pprof 分析
- 尝试将前端页面接入 Go 网关/认证服务
- 对比 Java：goroutine vs ThreadPool、channel vs BlockingQueue、net/http vs Spring MVC

## 学习建议
- 每章配合 `go doc` / `go help` 自查
- 优先使用标准库；三方库仅在必要时引入
- 保持小步提交、运行 `go test ./...`，并记录问题于 `plan.md`
