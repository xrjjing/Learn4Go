# RabbitMQ Demo（独立实战）

- 场景：订单创建异步处理 + 延迟关闭/死信。
- 组件：Go 后端生产/消费、前端演示页、Mock 数据、Docker Compose。
- 文档：`rabbitmq-demo/docs/guide.md`、`rabbitmq-demo/docs/api.md`、`rabbitmq-demo/docs/mock.md`。
- 一键启动（可选）：`docker compose -f deployments/docker-compose.rabbit.yml up -d`，浏览器访问 http://localhost:8088。已有集群可跳过 compose，设置 `RABBITMQ_URL=amqp://admin:1@localhost:5672/` 直接运行 `go run rabbitmq-demo/cmd/demoapp/main.go`。

目录一览：
```
rabbitmq-demo/
  cmd/demoapp/           # 后端服务（生产+消费+HTTP API+静态资源）
  internal/rabbit/       # RabbitMQ 封装
  internal/logstore/     # 内存日志存储
  web/rabbitmq-demo/     # 前端页面
  mock/messages.json     # 演示数据
  docs/                  # 使用指南/API/Mock 说明
  Dockerfile             # Demo 服务镜像构建
```
