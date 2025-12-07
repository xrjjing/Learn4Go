package rabbit

import "context"

// MQ 定义生产与消费的通用接口，便于注入假实现进行本地/测试。
type MQ interface {
	DeclareTopology(ctx context.Context) error
	Publish(ctx context.Context, msg DemoMessage) error
	Consume(ctx context.Context, queue string, handler func(DemoMessage) error) error
	Stats(ctx context.Context, queue string) (messages int, consumers int, err error)
	Close()
}
