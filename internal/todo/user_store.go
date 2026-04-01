package todo

// 本文件描述“用户数据从哪里来”。
//
// 当前仓库里，登录、当前用户、后台用户管理等功能都默认依赖 MemoryUserStore。
// 如果后续要接数据库用户体系，通常会先从这个接口层扩展，而不是直接改 handler。
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

// User 是认证和 RBAC 共同依赖的用户模型。
// PasswordHash 不对外输出，Role 会被 handler.go 用于权限分支。
type User struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserStore 把“用户来源”从业务层里抽象出来，便于后续替换为数据库或外部身份源。
type UserStore interface {
	Create(ctx context.Context, email, passwordHash string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id uint) (User, error)
}

// MemoryUserStore 是当前默认实现。
// todo.NewServer() 不额外注入 UserStore 时，就会使用它。
type MemoryUserStore struct {
	mu        sync.Mutex
	seq       uint
	users     map[string]User // email -> User
	usersByID map[uint]User   // id -> User
}

// NewMemoryUserStore 会预置 3 个演示用户，因此前端登录页可以开箱即用。
// NewMemoryUserStore：启动时预置三组演示账号，portal.html 与 todo-login.html 的默认凭据都来自这里。
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

// Create 用于注册或后台创建用户。当前实现会默认给新用户分配普通用户角色。
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

// FindByEmail 是 handleLogin 的关键依赖。登录失败时，通常先确认这里是否查到了正确用户。
// FindByEmail：登录流程的关键入口。
func (m *MemoryUserStore) FindByEmail(ctx context.Context, email string) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[email]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

// FindByID 会被 /v1/me、RBAC 和 refresh token 校验频繁调用。
// FindByID：认证中间件放行后，/v1/me 和 RBAC 都会继续从这里补全用户信息。
func (m *MemoryUserStore) FindByID(ctx context.Context, id uint) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.usersByID[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
