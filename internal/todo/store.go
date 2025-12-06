package todo

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
	"sync/atomic"
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
	ListPaged(page, pageSize int) ([]Todo, int, error) // 分页查询，返回数据和总数
	Create(title string, userID uint) (Todo, error)
	Get(id int) (Todo, bool, error)
	Toggle(id int, done bool) (Todo, bool, error)
	Delete(id int) (bool, error)
}

// Store 提供并发安全的内存存储 (实现 TodoStore 接口)
type Store struct {
	mu     sync.Mutex
	items  map[int]Todo
	nextID atomic.Int64 // 原子递增的 ID 生成器，确保并发安全
}

// NewStore 创建空存储。
func NewStore() *Store {
	s := &Store{items: make(map[int]Todo)}
	// 使用加密随机数初始化 ID 起始值，避免 ID 可预测
	var seed int64
	if err := binary.Read(rand.Reader, binary.BigEndian, &seed); err != nil {
		seed = time.Now().UnixNano()
	}
	s.nextID.Store(seed & 0x7FFFFFFF) // 确保为正数
	return s
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

// ListPaged 分页返回待办列表
// page 从 1 开始，pageSize 为每页条数
// 返回当前页数据和总条数
func (s *Store) ListPaged(page, pageSize int) ([]Todo, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 收集所有项目
	all := make([]Todo, 0, len(s.items))
	for _, v := range s.items {
		all = append(all, v)
	}
	total := len(all)

	// 计算偏移
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	if offset >= total {
		return []Todo{}, total, nil
	}

	end := offset + pageSize
	if end > total {
		end = total
	}

	return all[offset:end], total, nil
}

// Create 新建待办。
// 使用原子递增 ID，确保并发安全且不冲突
func (s *Store) Create(title string, userID uint) (Todo, error) {
	id := int(s.nextID.Add(1))

	t := Todo{
		ID:        id,
		UserID:    userID,
		Title:     title,
		CreatedAt: time.Now(),
	}

	s.mu.Lock()
	s.items[id] = t
	s.mu.Unlock()

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
