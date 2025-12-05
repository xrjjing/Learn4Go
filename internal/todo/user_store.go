package todo

import (
	"context"
	"errors"
	"sync"
	"time"
)

// 用户相关错误
var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
)

// User 用户实体
type User struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserStore 抽象用户存储。
type UserStore interface {
	Create(ctx context.Context, email, passwordHash string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id uint) (User, error)
}

// MemoryUserStore 内存用户存储（用于测试/演示）
type MemoryUserStore struct {
	mu        sync.Mutex
	seq       uint
	users     map[string]User // email -> User
	usersByID map[uint]User   // id -> User
}

// NewMemoryUserStore 创建内存存储并添加mock数据。
func NewMemoryUserStore() *MemoryUserStore {
	store := &MemoryUserStore{
		users:     make(map[string]User),
		usersByID: make(map[uint]User),
	}

	// 添加默认的mock数据（密码已预先使用bcrypt加密）
	// 原始密码：admin123, user123, demo123
	mockUsers := []struct {
		email string
		hash  string
		role  Role
	}{
		// admin123 的 bcrypt hash
		{"admin@example.com", "$2a$10$y3vLvny/mpVUAn.SUr6ZZOGQ7w82eUhjqVllRvGyqOdfEo55Q5F3O", RoleAdmin},
		// user123 的 bcrypt hash
		{"user@example.com", "$2a$10$3eH8RR2IAKBLx9NwxupasuydQuDCDZ.b1j5.LuCUbabM0aYBR11V6", RoleUser},
		// demo123 的 bcrypt hash
		{"demo@example.com", "$2a$10$XRWfjsbbS7AsARoLZlL6XONoZuEdCGRM0NMWKVFCDDGinD4hWRd0y", RoleGuest},
	}

	for _, u := range mockUsers {
		store.seq++
		user := User{
			ID:           store.seq,
			Email:        u.email,
			PasswordHash: u.hash,
			Role:         u.role,
			CreatedAt:    time.Now(),
		}
		store.users[u.email] = user
		store.usersByID[user.ID] = user
	}

	return store
}

func (m *MemoryUserStore) Create(ctx context.Context, email, passwordHash string) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[email]; ok {
		return User{}, ErrEmailExists
	}
	m.seq++
	u := User{
		ID:           m.seq,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         RoleUser, // 默认注册为普通用户
		CreatedAt:    time.Now(),
	}
	m.users[email] = u
	m.usersByID[u.ID] = u
	return u, nil
}

func (m *MemoryUserStore) FindByEmail(ctx context.Context, email string) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[email]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (m *MemoryUserStore) FindByID(ctx context.Context, id uint) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.usersByID[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
