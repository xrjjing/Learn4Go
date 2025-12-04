# 项目实战（从易到难）

| 编号 | 入口 | 目标 | 关键点 |
| --- | --- | --- | --- |
| 01 | `cmd/batchrename` | 命令行批量重命名（默认 dry-run，可 `--apply`） | flag、filepath、错误处理 |
| 03 | `cmd/todoapi` | 内存版 TODO REST API（标准库 net/http） | 路由、JSON、并发安全、日志中间件 |
| 05 | `cmd/logzip` | 将若干日志写入 ZIP，演示归档与校验 | archive/zip、io、校验和 |

> 编号与 Python 仓库保持节奏一致（缺口编号为预留后续扩展：日志分析、爬虫、worker、gateway 等）。
