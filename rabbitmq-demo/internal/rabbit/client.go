package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config 描述 RabbitMQ 拓扑与连接配置。
type Config struct {
	URL           string
	Exchange      string
	DelayExchange string
	DLXExchange   string
	WorkQueue     string
	DelayQueue    string
	DLXQueue      string
	RoutingWork   string
	RoutingDLX    string
	Prefetch      int
}

// DefaultConfig 返回默认演示配置。
func DefaultConfig() Config {
	return Config{
		URL:           "amqp://guest:guest@localhost:5672/",
		Exchange:      "orders.exchange",
		DelayExchange: "orders.delay.exchange",
		DLXExchange:   "orders.dlx.exchange",
		WorkQueue:     "orders.work.q",
		DelayQueue:    "orders.delay.q",
		DLXQueue:      "orders.dlx.q",
		RoutingWork:   "order.*",
		RoutingDLX:    "orders.dlx",
		Prefetch:      10,
	}
}

// DemoMessage 为演示用消息结构。
type DemoMessage struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	TTLMS   int64           `json:"ttl_ms,omitempty"`
}

// Client 封装连接与常用操作。
// Client 真实 RabbitMQ 客户端实现 MQ 接口。
type Client struct {
	conf   Config
	conn   *amqp.Connection
	ch     *amqp.Channel
	closed chan struct{}
}

// New 创建客户端并建立连接。
func New(cfg Config) (*Client, error) {
	if cfg.URL == "" {
		return nil, errors.New("缺少 RabbitMQ URL 配置")
	}
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("连接 RabbitMQ 失败: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建 channel 失败: %w", err)
	}
	if cfg.Prefetch <= 0 {
		cfg.Prefetch = 10
	}
	if err := ch.Qos(cfg.Prefetch, 0, false); err != nil {
		conn.Close()
		return nil, fmt.Errorf("设置 QoS 失败: %w", err)
	}
	return &Client{conf: cfg, conn: conn, ch: ch, closed: make(chan struct{})}, nil
}

// Close 关闭连接。
func (c *Client) Close() {
	select {
	case <-c.closed:
		return
	default:
		close(c.closed)
		_ = c.ch.Close()
		_ = c.conn.Close()
	}
}

// DeclareTopology 创建交换机与队列绑定。
func (c *Client) DeclareTopology(ctx context.Context) error {
	// 主交换机 topic
	if err := c.ch.ExchangeDeclare(c.conf.Exchange, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("声明主交换机失败: %w", err)
	}
	// 延迟交换机使用 direct
	if err := c.ch.ExchangeDeclare(c.conf.DelayExchange, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("声明延迟交换机失败: %w", err)
	}
	// DLX
	if err := c.ch.ExchangeDeclare(c.conf.DLXExchange, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("声明死信交换机失败: %w", err)
	}

	// 工作队列，绑定到主交换机
	if _, err := c.ch.QueueDeclare(c.conf.WorkQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    c.conf.DLXExchange,
		"x-dead-letter-routing-key": c.conf.RoutingDLX,
	}); err != nil {
		return fmt.Errorf("声明工作队列失败: %w", err)
	}
	if err := c.ch.QueueBind(c.conf.WorkQueue, c.conf.RoutingWork, c.conf.Exchange, false, nil); err != nil {
		return fmt.Errorf("绑定工作队列失败: %w", err)
	}

	// 延迟队列，TTL 后进入 DLX
	if _, err := c.ch.QueueDeclare(c.conf.DelayQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    c.conf.DLXExchange,
		"x-dead-letter-routing-key": c.conf.RoutingDLX,
	}); err != nil {
		return fmt.Errorf("声明延迟队列失败: %w", err)
	}
	if err := c.ch.QueueBind(c.conf.DelayQueue, c.conf.RoutingDLX, c.conf.DelayExchange, false, nil); err != nil {
		return fmt.Errorf("绑定延迟队列失败: %w", err)
	}

	// 死信队列
	if _, err := c.ch.QueueDeclare(c.conf.DLXQueue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("声明死信队列失败: %w", err)
	}
	if err := c.ch.QueueBind(c.conf.DLXQueue, c.conf.RoutingDLX, c.conf.DLXExchange, false, nil); err != nil {
		return fmt.Errorf("绑定死信队列失败: %w", err)
	}
	return nil
}

// Publish 发送消息，ttl>0 时进入延迟队列，否则进入主交换机。
func (c *Client) Publish(ctx context.Context, msg DemoMessage) error {
	if msg.ID == "" {
		return errors.New("消息缺少 id")
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	publish := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
		Timestamp:    time.Now(),
		MessageId:    msg.ID,
		Type:         msg.Type,
	}

	if msg.TTLMS > 0 {
		publish.Expiration = fmt.Sprintf("%d", msg.TTLMS)
		return c.ch.PublishWithContext(ctx, c.conf.DelayExchange, c.conf.RoutingDLX, false, false, publish)
	}
	return c.ch.PublishWithContext(ctx, c.conf.Exchange, msg.Type, false, false, publish)
}

// Consume 启动消费协程。
// handler 返回 error 时将 nack 且不重回队列，让消息进入 DLX。
func (c *Client) Consume(ctx context.Context, queue string, handler func(DemoMessage) error) error {
	deliveries, err := c.ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("注册消费者失败: %w", err)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-deliveries:
				if !ok {
					return
				}
				var msg DemoMessage
				if err := json.Unmarshal(d.Body, &msg); err != nil {
					_ = d.Nack(false, false)
					continue
				}
				if err := handler(msg); err != nil {
					_ = d.Nack(false, false)
					continue
				}
				_ = d.Ack(false)
			}
		}
	}()
	return nil
}

// Stats 返回队列的消息数与消费者数量。
func (c *Client) Stats(ctx context.Context, queue string) (int, int, error) {
	info, err := c.ch.QueueInspect(queue)
	if err != nil {
		return 0, 0, err
	}
	return info.Messages, info.Consumers, nil
}
