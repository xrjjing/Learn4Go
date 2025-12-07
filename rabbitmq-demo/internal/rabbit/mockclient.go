package rabbit

import (
	"context"
	"sync"
	"time"
)

// MockClient 提供内存版 MQ，便于无 RabbitMQ 时的演示与测试。
type MockClient struct {
	conf   Config
	mu     sync.Mutex
	queues map[string]chan DemoMessage
	closed chan struct{}
}

// NewMock 创建内存 MQ。
func NewMock(conf Config) *MockClient {
	return &MockClient{
		conf:   conf,
		queues: make(map[string]chan DemoMessage),
		closed: make(chan struct{}),
	}
}

func (m *MockClient) getQueue(name string) chan DemoMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	q, ok := m.queues[name]
	if !ok {
		q = make(chan DemoMessage, 256)
		m.queues[name] = q
	}
	return q
}

func (m *MockClient) DeclareTopology(ctx context.Context) error {
	// 内存队列按需创建，无需实际拓扑操作。
	_ = m.getQueue(m.conf.WorkQueue)
	_ = m.getQueue(m.conf.DelayQueue)
	_ = m.getQueue(m.conf.DLXQueue)
	return nil
}

func (m *MockClient) Publish(ctx context.Context, msg DemoMessage) error {
	if msg.TTLMS > 0 {
		// 延迟模拟：到期后进入 DLX 队列
		go func() {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(msg.TTLMS) * time.Millisecond):
				m.getQueue(m.conf.DLXQueue) <- msg
			}
		}()
		return nil
	}
	// 普通消息进入工作队列
	m.getQueue(m.conf.WorkQueue) <- msg
	return nil
}

func (m *MockClient) Consume(ctx context.Context, queue string, handler func(DemoMessage) error) error {
	q := m.getQueue(queue)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-q:
				if !ok {
					return
				}
				if err := handler(msg); err != nil {
					// 模拟进入死信队列
					m.getQueue(m.conf.DLXQueue) <- msg
				}
			}
		}
	}()
	return nil
}

func (m *MockClient) Close() {
	select {
	case <-m.closed:
		return
	default:
		close(m.closed)
		m.mu.Lock()
		defer m.mu.Unlock()
		for _, q := range m.queues {
			close(q)
		}
	}
}

// Stats 返回内存队列长度，消费者数量用 0 近似（无真实连接信息）。
func (m *MockClient) Stats(ctx context.Context, queue string) (int, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	q, ok := m.queues[queue]
	if !ok {
		return 0, 0, nil
	}
	return len(q), 0, nil
}
