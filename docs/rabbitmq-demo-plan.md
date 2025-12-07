# RabbitMQ 实战演练项目计划

## 1. 项目概览
- 目的：提供一个可运行的 RabbitMQ 生产/消费示例，配套前端演示页、使用指南、Mock 数据，使代码、文档、页面三者联动。
- 场景设定：以“订单创建→异步扣库存→延迟关闭未支付订单”为主线，同时演示普通队列、延迟/死信处理。
- 交付形式：同仓代码 + 前端演示页 + 文档（指南/API/Mock）+ 一键启动脚本。

## 2. 成功标准
- 一键启动：`docker compose up -d` 后，后端与 RabbitMQ 可用，前端可访问。
- 核心链路可视：前端可发消息，查看消费日志、死信/延迟队列状态；后端日志能看到发布确认与消费 ack。
- 文档齐备：README 顶部可跳转到使用指南、API 文档、Mock 说明，步骤可复现。
- 测试通过：基础单元/集成测试可在 60s 内完成；至少覆盖连接抽象、发布确认、重试/死信逻辑。

## 3. 功能范围
- 生产者：支持 direct/topic 发布，开启发布确认，读取 mock JSON 批量发送。
- 消费者：手动 ack，prefetch 可配，并发 worker；重试策略（重回队列/进入死信）；延迟队列示例。
- 前端演示：发送表单（JSON/文本），实时消费日志（SSE/WebSocket/轮询），队列堆积与死信数展示。
- 观测：暴露基础指标或调用管理端 API 获取队列长度、消费者数。

## 4. 技术选型
- 后端：Go + `github.com/rabbitmq/amqp091-go`；日志 `zap` 或内置 `slog`；配置 `viper` 可选。
- 前端：React + Vite（现有栈可复用）；UI 简洁，支持移动端。
- 基础设施：Docker Compose 启 RabbitMQ（含 management，必要时启延迟插件）、后端服务、前端容器；`.env` 统一配置。

## 5. 系统设计（简要）
- 交换机/队列
  - `orders.exchange` (topic)：路由 `order.created`、`order.cancelled`
  - 工作队列：`orders.work.q` → 正常消费
  - 延迟/死信：`orders.delay.exchange` + `orders.delay.q`（TTL 后路由到死信）
  - 死信队列：`orders.dlx.q`（存放失败/超时消息）
- 流程
  1. 前端提交订单消息 → 后端生产者发布到 topic 交换机。
  2. 消费者从工作队列获取，业务处理失败时可选择重试或 nack 进入死信。
  3. 延迟关闭：消息先进入延迟队列（TTL），到期转死信队列，由关闭订单消费者处理。

## 6. 目录结构建议
```
docs/
  rabbitmq-demo-plan.md      # 本计划
  rabbitmq-guide.md          # 使用指南（启动、调试、常见问题）
  rabbitmq-api.md            # 后端 API 说明（发送/查询日志接口）
  rabbitmq-mock.md           # Mock 数据格式与示例
cmd/
  rabbitmq-producer/         # 生产者 CLI
  rabbitmq-consumer/         # 消费者服务
internal/rabbitmq/           # 连接、发布、消费抽象（复用）
web/rabbitmq-demo/           # 前端演示页面
mock/messages.json           # 演示用消息
deployments/docker-compose.rabbit.yml # 专用 compose
```

## 7. 开发里程碑（建议节奏）
1) 环境与配置：完成 compose、创建 vhost/用户/权限，写 `.env.example`。  
2) 公共封装：`internal/rabbitmq` 封装连接、channel、发布确认、手动 ack、QoS。  
3) 最小链路：实现 `producer`、`consumer`，串通普通发布/消费。  
4) 延迟/死信：补 TTL+DLX 流程与消费端重试策略。  
5) 前端骨架：页面布局、发送表单、日志面板、队列状态区。  
6) 文档联动：补 `guide/api/mock` 文档，README 挂链接，脚本 `make demo`。  
7) 测试与验证：单元/集成测试 + 前端关键交互测试；本地跑通一键演示。  
8) 优化与验收：调优 prefetch/并发，补告警阈值建议，打标签/版本。

## 8. 文档与使用指南骨架
- 快速开始：前置依赖、`docker compose -f deployments/docker-compose.rabbit.yml up -d`、访问 URLs。
- 配置说明：连接串、交换机/队列/路由键、TTL/DLX 参数表。
- 操作步骤：发送消息、查看消费日志、模拟失败进入死信、触发延迟关闭。
- 故障排查：连接失败、堆积、消息丢失的常见原因与解决步骤。

## 9. Mock 数据规划
- 文件：`mock/messages.json`
- 字段示例：`id`（唯一键）、`type`（order.created/cancelled）、`payload`（订单摘要）、`ttl_ms`（可选延迟）。
- 用途：生产者 CLI 批量发送；前端显示模板；文档中提供示例。

## 10. 前端页面要点
- 区域：配置概览、发送面板、消费/死信日志、队列状态卡片。
- 交互：发送后返回 msgId；日志区自动滚动；失败条目高亮；队列长度定时刷新。
- 技术：调用后端 REST/SSE；轻量样式，保持 80 列可读。

## 11. 测试与验证
- 单测：连接管理、发布确认重试、消息序列化/反序列化。
- 集成：启动本地 RabbitMQ，验证发布→消费→死信/延迟链路。
- 前端：组件测试（表单校验、日志渲染），一次端到端流程。

## 12. 风险与待办
- 延迟插件许可与体积：若不启插件，采用 TTL+DLX 方案替代。
- 观测数据源：管理端 API 需要认证，可用后端代理抽象。
- CI 时间：确保测试 <60s；必要时分层运行。

---
后续执行顺序建议：先完成里程碑 1-3（后端最小链路），再并行推进前端骨架与文档，最后补测试与优化。 
