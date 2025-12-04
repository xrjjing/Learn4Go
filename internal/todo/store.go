package todo

import (
	"math/rand"
	"sync"
	"time"
)

// Todo 表示单条待办。
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// Store 提供并发安全的内存存储。
type Store struct {
	mu    sync.Mutex
	items map[int]Todo
}

// NewStore 创建空存储。
func NewStore() *Store {
	return &Store{items: make(map[int]Todo)}
}

// List 返回全部待办的拷贝。
func (s *Store) List() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Todo, 0, len(s.items))
	for _, v := range s.items {
		out = append(out, v)
	}
	return out
}

// Create 新建待办。
func (s *Store) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := rand.Intn(1_000_000)
	t := Todo{ID: id, Title: title, CreatedAt: time.Now()}
	s.items[id] = t
	return t
}

// Toggle 设置完成状态。
func (s *Store) Toggle(id int, done bool) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.items[id]
	if !ok {
		return Todo{}, false
	}
	t.Done = done
	s.items[id] = t
	return t, true
}

// Delete 删除待办。
func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}
