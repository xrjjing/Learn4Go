# RabbitMQ Demo 使用指南

## 背景

演示“订单创建 → 异步处理 → 延迟关闭/死信”完整链路，配套前端页面、Mock 数据和可运行的 Go 代码。

## 使用方式对比

### 方式 A：自带 Docker Compose（一键体验）

- 适合：本机无 RabbitMQ，想快速跑通。
- 启动：
  ```bash
  docker compose -f deployments/docker-compose.rabbit.yml up -d
  ```
  - 管理端：http://localhost:15672 （guest/guest）
  - Demo 服务：http://localhost:8088
- 访问前端：打开 `http://localhost:8088/`，发送消息、查看日志。

### 方式 B：接入你已有的双节点集群

- 适合：已运行 rabbitmq-node-01/02（3.6.12），端口 5672/5673，账号 `admin/1`，vhost `/`。
- 运行 Demo（跳过 compose）：
  ```bash
  RABBITMQ_URL=amqp://admin:1@localhost:5672/ PORT=8088 \
  go run rabbitmq-demo/cmd/demoapp/main.go
  ```
  若连节点 2，则改为 5673。前端同样访问 http://localhost:8088/。

### 方式 C：Fake 模式（无 RabbitMQ/CI）

- 适合：离线演示、前后端联调或 CI。
  ```bash
  RABBITMQ_FAKE=1 PORT=8088 go run rabbitmq-demo/cmd/demoapp/main.go
  # 或仅跑测试
  GOCACHE=$(pwd)/.gocache RABBITMQ_FAKE=1 go test ./rabbitmq-demo/cmd/demoapp -run Test -count=1
  ```
- 前端仍访问 http://localhost:8088/，消息在内存模拟。

---

## 常规操作

- 发送单条：前端表单或 `POST /api/messages`
- 批量 Mock：前端按钮或 `POST /api/messages/batch`
- 查看日志：前端“消费与死信日志”或 `GET /api/logs`

## 配置项

通过环境变量覆盖默认值：

- `RABBITMQ_URL` (默认 `amqp://guest:guest@localhost:5672/`)
- `RABBITMQ_EXCHANGE` / `RABBITMQ_DELAY_EXCHANGE` / `RABBITMQ_DLX_EXCHANGE`
- `RABBITMQ_WORK_QUEUE` / `RABBITMQ_DELAY_QUEUE` / `RABBITMQ_DLX_QUEUE`
- `RABBITMQ_ROUTING_WORK` / `RABBITMQ_ROUTING_DLX`
- `RABBITMQ_PREFETCH` (默认 10)
- `LOG_CAPACITY` (默认 400，范围 1-10000) - 内存日志存储容量，超出范围将回退到默认值
- `PORT` (默认 8088)
- `RABBITMQ_FAKE`=1 启用内存模式

## 高可用/多节点补充

- 你的双节点（5672/5673）可作为主/备入口，客户端无需同时连两端；生产环境建议使用 TCP 负载均衡或 DNS RR。
- 若要队列级高可用，可在 3.6.12 使用镜像队列（HA policies）；若升级到 3.8+ 推荐仲裁队列（quorum queues）。
- 管理端与指标：当前前端“运行状态”依赖 `/api/status`，通过 AMQP `QueueInspect` 获取消息/消费者数；如需更全面指标，可在后端代理管理 API 或 Prometheus exporter。

## 功能验证路径

1. **正常消费**：发送 `order.created`，应在日志看到 send → consume。
2. **业务失败进入死信**：发送 `order.fail`，消费端会返回错误，消息进入 DLX，日志出现 `dlx`。
3. **延迟关闭**：发送 `order.closed` 且设置 `ttl_ms`（例如 15000），到期后出现在死信队列，由关闭订单消费者处理。

## 常见问题

- 连接失败：确认 RabbitMQ 端口开放，URL 用户/密码正确。
- 消息不进入死信：检查队列参数 `x-dead-letter-exchange`、`x-dead-letter-routing-key` 是否与配置一致，消费者是否 `nack` 且不重回队列。
- 页面无日志：确认 Demo 服务端口 8088 可访问；查看后台日志是否有报错。

## 目录快速索引

- 代码：`rabbitmq-demo/cmd/demoapp`、`rabbitmq-demo/internal/rabbit`
- 前端：`rabbitmq-demo/web/rabbitmq-demo`
- Mock：`rabbitmq-demo/mock/messages.json`
- 文档：`rabbitmq-demo/docs/*`
