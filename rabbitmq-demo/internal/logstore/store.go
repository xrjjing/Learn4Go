package logstore

import (
	"sync"
	"time"
)

// Entry 表示一次事件日志。
type Entry struct {
	Time    time.Time `json:"time"`
	Kind    string    `json:"kind"` // send/consume/dlx/error
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
}

// Store 提供固定容量的线程安全日志存储。
type Store struct {
	mu       sync.Mutex
	capacity int
	entries  []Entry
}

// New 创建固定容量的日志存储。
func New(capacity int) *Store {
	if capacity <= 0 {
		capacity = 200
	}
	return &Store{capacity: capacity}
}

// Add 追加一条日志，超过容量则丢弃最早的记录。
func (s *Store) Add(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.entries) >= s.capacity {
		drop := len(s.entries) - s.capacity + 1
		s.entries = s.entries[drop:]
	}
	s.entries = append(s.entries, e)
}

// List 返回当前所有日志的拷贝，按时间顺序。
func (s *Store) List() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}
