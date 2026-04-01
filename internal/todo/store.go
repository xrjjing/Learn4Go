package todo

// 本文件是 TODO 数据层的内存实现。
//
// 适用场景：
// - 本地学习
// - 不依赖数据库时的快速联调
// - `TODO_STORAGE=memory` 模式
//
// 上游调用方主要是 handler.go 中的 TODO 路由；
// 若你想确认“数据到底有没有真正写进去”，可以直接从 Create/Get/Toggle/Delete 看。
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

// TodoStore 是 handler 层依赖的抽象边界。
//
// 这意味着 handler 不关心底层是内存 map 还是数据库，只关心接口语义是否一致。
// 将错误显式向上返回，可以让 API 层正确区分 404、500 和真正的业务失败。
type TodoStore interface {
	List() ([]Todo, error)
	ListByUser(userID uint) ([]Todo, error)
	ListPaged(page, pageSize int) ([]Todo, int, error) // 分页查询，返回数据和总数
	Create(title string, userID uint) (Todo, error)
	Get(id int) (Todo, bool, error)
	Toggle(id int, done bool) (Todo, bool, error)
	Delete(id int) (bool, error)
}

// Store 是 TodoStore 的内存版实现。
// 使用 mutex 保护 items，使用 atomic.Int64 生成 ID，便于在并发场景下保持简单可读。
type Store struct {
	mu     sync.Mutex
	items  map[int]Todo
	nextID atomic.Int64 // 原子递增的 ID 生成器，确保并发安全
}

// NewStore 会初始化空 map，并给 nextID 一个随机起点，避免演示环境中 ID 过于可预测。
// NewStore：创建并初始化并发安全的内存仓库。
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

// List 给管理员/访客查看全量 TODO 使用。普通用户通常走 ListByUser。
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
// ListByUser：普通用户查询自己的待办时会落到这里。
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

// Create 是最常见的写入入口。
// 调用链通常是：POST /v1/todos → handler.go → Store.Create。
// Create：POST /v1/todos 的最终写入点之一。
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

// Toggle 由 PUT /v1/todos/{id} 调用，只负责更新 done，不修改其他字段。
// Toggle：PUT /v1/todos/{id} 的最终更新点之一。
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

// Delete 由 DELETE /v1/todos/{id} 调用。若返回 false，handler 会翻译成 404。
// Delete：DELETE /v1/todos/{id} 的最终删除点之一。
func (s *Store) Delete(id int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false, nil
	}
	delete(s.items, id)
	return true, nil
}
