# RabbitMQ Demo Mock 数据

文件：`rabbitmq-demo/mock/messages.json`

字段说明：
- `id`：消息唯一标识（可为空，服务端会自动生成）。
- `type`：路由键，对应交换机绑定（如 `order.created`）。
- `payload`：任意 JSON 业务字段。
- `ttl_ms`：可选，毫秒级 TTL；>0 时消息被发送到延迟队列，到期进入死信队列。

当前示例：
1. `order.created` 正常消费。
2. `order.cancelled` 正常消费。
3. `order.fail` 消费端返回错误，消息进入死信。
4. `order.closed` 设置 15s TTL，到期后进入死信队列，模拟超时关闭订单。
