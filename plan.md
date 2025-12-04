# Learn4Go 迭代计划

## 目标
- 按 Python learn 仓库的节奏与难度递进，构建 Go 版“语言基础 + 工程化/框架 + 项目实战 + 前后端联动”练习场。
- 前端沿用原 UI/UX，后端与文档聚焦 Go 示例与最佳实践。

## 当前里程碑
- ✅ 仓库 scaffold：基础目录、README、MIT License、前端静态页复制（现位于 `web/`）
- ✅ Go 标准布局：`cmd/*` 可执行入口、示例迁移完成；业务逻辑抽到 `internal/*`
- ✅ Makefile：fmt / vet / test / lint 占位
- ⏳ 文档初稿：Go 学习规划、Java/Go 对照表、架构草图、快速开始
- ⏳ 基础章节：10 个核心语法与并发章节草稿
- ✅ 实战示例：3 个可运行示例（命令行批量重命名、HTTP mini TODO、ZIP 归档），均可 `go run ./cmd/...`

## 任务清单
1. 文档
   - [ ] 完成 `docs/Go学习规划_Java开发者版.md` 细化练习清单
   - [ ] 完成 `docs/Java_vs_Go_CheatSheet.md` 并补充代码对照
   - [ ] 将菜鸟教程风格的示例整理入 `01.Go语言基础`（需联网检索确认细节）
   - [ ] 前端接口与后端映射文档（web/ 与 cmd 服务对应表）
2. 基础章节
   - [ ] 为每章补充最小可运行代码片段与练习题（可放 `examples/`）
   - [ ] 添加 `go test` 示例与基准测试示例（部分已在 todo/batchrename/logzip）
3. 项目实战
   - [ ] 01_cli_batch_rename：默认 dry-run，可加 `--apply`（入口 `cmd/batchrename`）
   - [ ] 03_http_todo: 使用标准库 `net/http` 的内存版 CRUD（入口 `cmd/todoapi`）
   - [ ] 05_log_archive_zip: 归档与校验示例（入口 `cmd/logzip`）
   - [ ] 后续：并发 worker、gRPC 小示例、网关转发/鉴权（对标 Python 网关）
4. 前端联动
   - [ ] 将前端接口路径替换为 Go 服务地址（静态资源现位于 `web/`）
   - [ ] 增加示例 token/网关调用说明
5. 质量与工具链
   - [ ] 增加 `Makefile`（lint/test/format），`golangci-lint` 配置
   - [ ] CI 占位（GitHub Actions: go test）

## 风险与依赖
- 需访问菜鸟教程、GitHub、Stack Overflow 获取最新示例与最佳实践（当前网络受限，待授权）。
- Gin/Fiber 等三方依赖需要 `go get`，如网络受限将使用标准库版本示例。

## 下一步
- 获得网络访问许可后，拉取菜鸟教程 Go 基础示例，充实章节代码与练习。
- 完成 3 个实战示例的可运行代码，并在 `plan.md` 更新状态。
