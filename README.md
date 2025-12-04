# Go 学习与实战仓库（对标 Python learn 项目）

本仓库是基于你原有的 Python 学习仓库 `/Users/xrj/PycharmProjects/learn` 的结构与节奏，面向 Go 语言的系统化练习场：**语言基础 → 工程化/框架 → 项目实战 → 前后端联动**。前端沿用原仓库的 UI/UX（已复制到 `web/`），文本与示例围绕 Go 场景更新。

## 📂 目录概览
- `docs/`：学习规划、Java 对照速查、架构与快速开始
- `01.Go语言基础/`：按章节的语法与并发特性笔记 + 代码片段
- `02.开发环境及框架介绍/`：Go 工程化、Gin/Echo/Fiber 对比、gRPC 与配置管理
- `03.项目实战/`：实战说明文档（代码放在 `cmd/` 下的独立可执行入口）
- `cmd/`：可执行入口（符合 Go 标准布局）
  - `cmd/learn4go/`：最小入口示例
  - `cmd/todoapi/`：内存版 TODO REST API（net/http）
  - `cmd/batchrename/`：批量重命名 CLI（默认 dry-run）
  - `cmd/logzip/`：日志 ZIP 归档示例
- `internal/`：复用/业务逻辑（todo、batchrename、logzip）
- `examples/`：基础章节的可运行小示例
- `web/`：沿用 Python 仓库的静态页面，用于演示网关/认证/日志分析的调用链（文本将逐步替换为 Go 版文案）

## 🚀 快速开始
```bash
go run ./cmd/learn4go
# TODO API
go run ./cmd/todoapi
# 批量重命名（默认 dry-run）
go run ./cmd/batchrename --dir=. --suffix=.txt --prefix=new_
# 日志归档 ZIP
go run ./cmd/logzip
# 质量检查
make fmt vet test
```

## 📚 推荐学习路径
1. 阅读 `docs/Go学习规划_Java开发者版.md`，了解阶段目标与检查清单。
2. 按顺序完成 `01.Go语言基础` 各章节的小练习，结合 `docs/Java_vs_Go_CheatSheet.md` 对照 Java 思维。
3. 进入 `02.开发环境及框架介绍`，搭建 Go 开发环境、模块管理与常见 Web 框架认知。
4. 按 `03.项目实战` 的编号从易到难实现：命令行 → 日志处理 → REST API → 并发 → 网关代理。

## 🔗 与 Python 仓库的映射
- 目录层级与编号保持一致，便于“对照迁移”与“跨语言对比”。
- 前端页面结构一致（移至 `web/`），后续将替换为 Go 后端的接口路径示例。
- 进阶专题（并发、网络、序列化、测试）在 Go 版放入基础章节与实战项目中逐步覆盖。

## 📄 许可证
采用 MIT License，详见 `LICENSE`。

## 📌 状态同步
- 迭代计划：见 `plan.md`
- 文档与练习：见各目录内 README / 章节文件
- 待办与下一步：`plan.md` 持续更新
- 前端映射：见 `docs/frontend_mapping.md`
