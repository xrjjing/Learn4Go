package logstore

import (
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// 测试正常容量
	s := New(100)
	if s.capacity != 100 {
		t.Errorf("容量不正确: got %d, want 100", s.capacity)
	}

	// 测试零容量应设为默认值
	s2 := New(0)
	if s2.capacity != 200 {
		t.Errorf("零容量应设为默认值 200: got %d", s2.capacity)
	}

	// 测试负数容量应设为默认值
	s3 := New(-10)
	if s3.capacity != 200 {
		t.Errorf("负数容量应设为默认值 200: got %d", s3.capacity)
	}
}

func TestStore_Add(t *testing.T) {
	s := New(3)

	e1 := Entry{Time: time.Now(), Kind: "send", ID: "1", Type: "order.created", Message: "msg1"}
	e2 := Entry{Time: time.Now(), Kind: "consume", ID: "2", Type: "order.created", Message: "msg2"}
	e3 := Entry{Time: time.Now(), Kind: "dlx", ID: "3", Type: "order.fail", Message: "msg3"}
	e4 := Entry{Time: time.Now(), Kind: "error", ID: "4", Type: "order.fail", Message: "msg4"}

	s.Add(e1)
	s.Add(e2)
	s.Add(e3)

	entries := s.List()
	if len(entries) != 3 {
		t.Errorf("日志数量不正确: got %d, want 3", len(entries))
	}

	// 超过容量，应丢弃最早的
	s.Add(e4)
	entries = s.List()
	if len(entries) != 3 {
		t.Errorf("日志数量应保持为容量上限: got %d, want 3", len(entries))
	}

	// 验证最早的 e1 被丢弃
	if entries[0].ID != "2" {
		t.Errorf("最早的日志应被丢弃: got ID %s, want 2", entries[0].ID)
	}
	if entries[2].ID != "4" {
		t.Errorf("最新的日志应该是 e4: got ID %s, want 4", entries[2].ID)
	}
}

func TestStore_List(t *testing.T) {
	s := New(10)

	// 空列表
	entries := s.List()
	if len(entries) != 0 {
		t.Errorf("初始列表应为空: got %d", len(entries))
	}

	// 添加日志
	e1 := Entry{Time: time.Now(), Kind: "send", ID: "1", Message: "test"}
	s.Add(e1)

	entries = s.List()
	if len(entries) != 1 {
		t.Errorf("日志数量不正确: got %d, want 1", len(entries))
	}

	// 验证返回的是拷贝，修改不影响原数据
	entries[0].Message = "modified"
	entries2 := s.List()
	if entries2[0].Message != "test" {
		t.Errorf("修改拷贝不应影响原数据: got %s, want test", entries2[0].Message)
	}
}

func TestStore_Concurrency(t *testing.T) {
	s := New(1000)
	var wg sync.WaitGroup

	// 并发添加
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				e := Entry{
					Time:    time.Now(),
					Kind:    "send",
					ID:      string(rune(id*10 + j)),
					Type:    "test",
					Message: "concurrent test",
				}
				s.Add(e)
			}
		}(i)
	}

	wg.Wait()

	entries := s.List()
	if len(entries) != 1000 {
		t.Errorf("并发添加后日志数量不正确: got %d, want 1000", len(entries))
	}
}

func TestEntry_Fields(t *testing.T) {
	now := time.Now()
	e := Entry{
		Time:    now,
		Kind:    "consume",
		ID:      "test-id",
		Type:    "order.created",
		Message: "test message",
	}

	if e.Time != now {
		t.Error("Time 字段不匹配")
	}
	if e.Kind != "consume" {
		t.Errorf("Kind 字段不匹配: got %s, want consume", e.Kind)
	}
	if e.ID != "test-id" {
		t.Errorf("ID 字段不匹配: got %s, want test-id", e.ID)
	}
	if e.Type != "order.created" {
		t.Errorf("Type 字段不匹配: got %s, want order.created", e.Type)
	}
	if e.Message != "test message" {
		t.Errorf("Message 字段不匹配: got %s, want test message", e.Message)
	}
}
