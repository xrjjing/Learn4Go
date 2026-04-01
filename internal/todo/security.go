package todo

// 本文件处理“JWT 之外”的安全细节。
//
// 重点包括：
// 1. 登录失败锁定策略
// 2. refresh token 的生成、存储、校验和旋转
// 3. 邮箱与密码强度校验
//
// 如果登录接口出现 429、refresh token 失效、密码格式被拒绝，优先看这里。
import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/mail"
	"time"
	"unicode"
)

const (
	maxLoginFailures   = 5                // 最大连续失败次数
	loginLockDuration  = 10 * time.Minute // 锁定时长
	loginFailureWindow = 15 * time.Minute // 统计失败窗口
	defaultRefreshTTL  = 7 * 24 * time.Hour
)

// loginFailure 用于记录同一邮箱在滑动时间窗内的失败次数和锁定状态。
type loginFailure struct {
	count        int
	lastFailedAt time.Time
	lockedUntil  time.Time
}

// refreshSession 表示 refresh token 在服务端内存表里的会话状态。
type refreshSession struct {
	userID    uint
	expiresAt time.Time
}

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshExpired      = errors.New("refresh token expired")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrWeakPassword        = errors.New("password must be at least 8 characters with letters and numbers")
)

// ValidateEmail 在注册和后台创建用户时做格式校验。
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword 目前要求长度不少于 8，且同时包含字母和数字。
// ValidatePassword：密码强度规则，避免弱口令直接写入用户仓库。
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	var hasLetter, hasDigit bool
	for _, c := range password {
		if unicode.IsLetter(c) {
			hasLetter = true
		}
		if unicode.IsDigit(c) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return ErrWeakPassword
	}
	return nil
}

// recordLoginFailure 是 handleLogin 的防暴力破解辅助函数。
// recordLoginFailure：登录失败后更新计数，并在必要时返回锁定剩余时长。
func (s *Server) recordLoginFailure(email string) (bool, time.Duration) {
	s.loginMu.Lock()
	defer s.loginMu.Unlock()

	now := time.Now()
	rec, ok := s.loginFailures[email]
	if !ok {
		rec = &loginFailure{}
		s.loginFailures[email] = rec
	}

	// 窗口外重置计数
	if now.Sub(rec.lastFailedAt) > loginFailureWindow {
		rec.count = 0
	}
	rec.count++
	rec.lastFailedAt = now

	if rec.count >= maxLoginFailures {
		rec.lockedUntil = now.Add(loginLockDuration)
	}

	if now.Before(rec.lockedUntil) {
		return true, rec.lockedUntil.Sub(now)
	}
	return false, 0
}

// clearLoginFailure 成功登录后清空失败计数
func (s *Server) clearLoginFailure(email string) {
	s.loginMu.Lock()
	defer s.loginMu.Unlock()
	delete(s.loginFailures, email)
}

// checkLock 判断是否处于锁定期
func (s *Server) checkLock(email string) (bool, time.Duration) {
	s.loginMu.Lock()
	defer s.loginMu.Unlock()
	rec, ok := s.loginFailures[email]
	if !ok || rec.lockedUntil.IsZero() {
		return false, 0
	}
	now := time.Now()
	if now.Before(rec.lockedUntil) {
		return true, rec.lockedUntil.Sub(now)
	}
	return false, 0
}

// generateRefreshToken 生成随机 refresh token，本身不包含用户信息；用户绑定关系存放在 refreshStore。
// generateRefreshToken：生成高熵随机串，避免 refresh token 可预测。
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// issueTokens 是登录成功和 refresh 成功后的共同出口。
//
// 调用链：handleLogin / handleRefresh -> issueTokens -> 返回 access + refresh。
// issueTokens：登录/刷新成功后的统一签发入口，同时旋转旧 refresh token。
func (s *Server) issueTokens(user User) (string, string, int, error) {
	access, err := s.jwtManager.Generate(user.ID)
	if err != nil {
		return "", "", 0, err
	}
	refresh, err := generateRefreshToken()
	if err != nil {
		return "", "", 0, err
	}

	s.refreshMu.Lock()
	// 清理该用户的旧 refresh，确保单一有效链
	for token, session := range s.refreshStore {
		if session.userID == user.ID {
			delete(s.refreshStore, token)
		}
	}
	s.refreshStore[refresh] = refreshSession{
		userID:    user.ID,
		expiresAt: time.Now().Add(s.refreshTTL),
	}
	s.refreshMu.Unlock()

	return access, refresh, int(s.jwtManager.cfg.TTL.Seconds()), nil
}

// validateRefresh 负责把 refresh token 重新映射回用户，并判断是否过期。
// validateRefresh：/v1/refresh 处理器进入业务前先调用这里确认 refresh token 是否有效。
func (s *Server) validateRefresh(token string) (User, error) {
	s.refreshMu.Lock()
	session, ok := s.refreshStore[token]
	if !ok {
		s.refreshMu.Unlock()
		return User{}, ErrInvalidRefreshToken
	}
	if time.Now().After(session.expiresAt) {
		delete(s.refreshStore, token)
		s.refreshMu.Unlock()
		return User{}, ErrRefreshExpired
	}
	s.refreshMu.Unlock()

	// 拉取用户信息
	user, err := s.userStore.FindByID(context.Background(), session.userID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// dropRefresh 在 refresh 轮换时使旧 token 立即失效，避免并存。
func (s *Server) dropRefresh(token string) {
	s.refreshMu.Lock()
	delete(s.refreshStore, token)
	s.refreshMu.Unlock()
}
