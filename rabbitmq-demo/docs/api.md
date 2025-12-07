# RabbitMQ Demo API 文档

基地址：`http://localhost:8088`

## POST /api/messages
- 功能：发送单条消息。
- 请求体示例
```json
{
  "type": "order.created",
  "ttl_ms": 15000,
  "payload": {"orderId": "A1001"}
}
```
- 响应
```json
{ "id": "生成的消息ID" }
```

## POST /api/messages/batch
- 功能：读取 `mock/messages.json` 批量发送。
- 请求体：无。
- 响应
```json
{ "published": 4 }
```

## GET /api/logs
- 功能：返回最近的事件日志（发送/消费/死信/错误）。
- 响应字段
  - `time`：时间戳
  - `kind`：send / consume / dlx / error
  - `id`：消息 ID
  - `type`：消息类型
  - `message`：描述

## GET /api/mock
- 功能：返回 Mock 消息数组，供前端展示或调试。

## GET /api/health
- 功能：健康检查，返回 `ok`。

## GET /api/status
- 功能：返回当前运行模式及队列指标。
- 响应示例
```json
{
  "mode": "real",
  "rabbit_url": "amqp://admin:1@localhost:5672/",
  "queues": [
    {"name":"orders.work.q","messages":0,"consumers":1},
    {"name":"orders.delay.q","messages":0,"consumers":0},
    {"name":"orders.dlx.q","messages":0,"consumers":1}
  ]
}
```
