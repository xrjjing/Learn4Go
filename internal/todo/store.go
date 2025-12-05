package todo

import (
	"math/rand"
	"sync"
	"time"
)

// Todo 表示单条待办。
type Todo struct {
	ID        int       `json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoStore 定义存储接口 (内存/数据库实现此接口)
// 将错误透传给上层，避免静默失败伪装成 2xx/404
type TodoStore interface {
	List() ([]Todo, error)
	ListByUser(userID uint) ([]Todo, error)
	Create(title string, userID uint) (Todo, error)
	Get(id int) (Todo, bool, error)
	Toggle(id int, done bool) (Todo, bool, error)
	Delete(id int) (bool, error)
}

// Store 提供并发安全的内存存储 (实现 TodoStore 接口)
type Store struct {
	mu    sync.Mutex
	items map[int]Todo
}

// NewStore 创建空存储。
func NewStore() *Store {
	return &Store{items: make(map[int]Todo)}
}

// List 返回全部待办的拷贝。
func (s *Store) List() ([]Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Todo, 0, len(s.items))
	for _, v := range s.items {
		out = append(out, v)
	}
	return out, nil
}

// ListByUser 返回指定用户的待办列表
func (s *Store) ListByUser(userID uint) ([]Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Todo, 0)
	for _, v := range s.items {
		if v.UserID == userID {
			out = append(out, v)
		}
	}
	return out, nil
}

// Create 新建待办。
func (s *Store) Create(title string, userID uint) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := rand.Intn(1_000_000)
	t := Todo{
		ID:        id,
		UserID:    userID,
		Title:     title,
		CreatedAt: time.Now(),
	}
	s.items[id] = t
	return t, nil
}

// Get 获取指定ID的待办
func (s *Store) Get(id int) (Todo, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.items[id]
	return t, ok, nil
}

// Toggle 设置完成状态。
func (s *Store) Toggle(id int, done bool) (Todo, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.items[id]
	if !ok {
		return Todo{}, false, nil
	}
	t.Done = done
	s.items[id] = t
	return t, true, nil
}

// Delete 删除待办。
func (s *Store) Delete(id int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false, nil
	}
	delete(s.items, id)
	return true, nil
}
