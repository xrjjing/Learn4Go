package rabbit

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestDemoMessageJSON(t *testing.T) {
	msg := DemoMessage{
		ID:      "test-001",
		Type:    "order.created",
		Payload: json.RawMessage(`{"orderId":"A001"}`),
		TTLMS:   5000,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var decoded DemoMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	if decoded.ID != msg.ID {
		t.Errorf("ID 不匹配: got %s, want %s", decoded.ID, msg.ID)
	}
	if decoded.Type != msg.Type {
		t.Errorf("Type 不匹配: got %s, want %s", decoded.Type, msg.Type)
	}
	if decoded.TTLMS != msg.TTLMS {
		t.Errorf("TTLMS 不匹配: got %d, want %d", decoded.TTLMS, msg.TTLMS)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.URL == "" {
		t.Error("URL 不应为空")
	}
	if cfg.Exchange == "" {
		t.Error("Exchange 不应为空")
	}
	if cfg.Prefetch <= 0 {
		t.Error("Prefetch 应为正数")
	}
}

func TestMockClient_Publish(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mock := NewMock(DefaultConfig())
	defer mock.Close()

	if err := mock.DeclareTopology(ctx); err != nil {
		t.Fatalf("声明拓扑失败: %v", err)
	}

	msg := DemoMessage{
		ID:      "mock-001",
		Type:    "order.created",
		Payload: json.RawMessage(`{"orderId":"M001"}`),
	}

	// 测试普通消息发布
	if err := mock.Publish(ctx, msg); err != nil {
		t.Fatalf("发布消息失败: %v", err)
	}

	// 验证消息进入工作队列
	workMessages, _, err := mock.Stats(ctx, mock.conf.WorkQueue)
	if err != nil {
		t.Fatalf("获取队列状态失败: %v", err)
	}
	if workMessages != 1 {
		t.Errorf("工作队列消息数不正确: got %d, want 1", workMessages)
	}
}

func TestMockClient_Consume(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mock := NewMock(DefaultConfig())
	defer mock.Close()

	if err := mock.DeclareTopology(ctx); err != nil {
		t.Fatalf("声明拓扑失败: %v", err)
	}

	received := make(chan DemoMessage, 1)
	handler := func(msg DemoMessage) error {
		received <- msg
		return nil
	}

	if err := mock.Consume(ctx, mock.conf.WorkQueue, handler); err != nil {
		t.Fatalf("启动消费者失败: %v", err)
	}

	msg := DemoMessage{
		ID:      "mock-002",
		Type:    "order.created",
		Payload: json.RawMessage(`{"orderId":"M002"}`),
	}

	if err := mock.Publish(ctx, msg); err != nil {
		t.Fatalf("发布消息失败: %v", err)
	}

	select {
	case recv := <-received:
		if recv.ID != msg.ID {
			t.Errorf("接收到的消息 ID 不匹配: got %s, want %s", recv.ID, msg.ID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("消费超时")
	}
}

func TestMockClient_DLX(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mock := NewMock(DefaultConfig())
	defer mock.Close()

	if err := mock.DeclareTopology(ctx); err != nil {
		t.Fatalf("声明拓扑失败: %v", err)
	}

	// 工作队列消费者返回错误，消息应进入死信队列
	workHandler := func(msg DemoMessage) error {
		return &struct{ error }{error: nil} // 模拟错误
	}

	dlxReceived := make(chan DemoMessage, 1)
	dlxHandler := func(msg DemoMessage) error {
		dlxReceived <- msg
		return nil
	}

	if err := mock.Consume(ctx, mock.conf.WorkQueue, workHandler); err != nil {
		t.Fatalf("启动工作队列消费者失败: %v", err)
	}

	if err := mock.Consume(ctx, mock.conf.DLXQueue, dlxHandler); err != nil {
		t.Fatalf("启动死信队列消费者失败: %v", err)
	}

	msg := DemoMessage{
		ID:      "mock-003",
		Type:    "order.fail",
		Payload: json.RawMessage(`{"orderId":"M003"}`),
	}

	if err := mock.Publish(ctx, msg); err != nil {
		t.Fatalf("发布消息失败: %v", err)
	}

	select {
	case recv := <-dlxReceived:
		if recv.ID != msg.ID {
			t.Errorf("死信队列接收到的消息 ID 不匹配: got %s, want %s", recv.ID, msg.ID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("死信队列消费超时")
	}
}

func TestMockClient_Delay(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mock := NewMock(DefaultConfig())
	defer mock.Close()

	if err := mock.DeclareTopology(ctx); err != nil {
		t.Fatalf("声明拓扑失败: %v", err)
	}

	dlxReceived := make(chan DemoMessage, 1)
	dlxHandler := func(msg DemoMessage) error {
		dlxReceived <- msg
		return nil
	}

	if err := mock.Consume(ctx, mock.conf.DLXQueue, dlxHandler); err != nil {
		t.Fatalf("启动死信队列消费者失败: %v", err)
	}

	msg := DemoMessage{
		ID:      "mock-004",
		Type:    "order.closed",
		Payload: json.RawMessage(`{"orderId":"M004"}`),
		TTLMS:   500, // 500ms 延迟
	}

	start := time.Now()
	if err := mock.Publish(ctx, msg); err != nil {
		t.Fatalf("发布延迟消息失败: %v", err)
	}

	select {
	case recv := <-dlxReceived:
		elapsed := time.Since(start)
		if elapsed < 400*time.Millisecond {
			t.Errorf("消息到达过早: %v", elapsed)
		}
		if recv.ID != msg.ID {
			t.Errorf("死信队列接收到的消息 ID 不匹配: got %s, want %s", recv.ID, msg.ID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("延迟消息未到达死信队列")
	}
}
