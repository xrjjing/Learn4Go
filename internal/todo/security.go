package todo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

const (
	maxLoginFailures   = 5                // 最大连续失败次数
	loginLockDuration  = 10 * time.Minute // 锁定时长
	loginFailureWindow = 15 * time.Minute // 统计失败窗口
	defaultRefreshTTL  = 7 * 24 * time.Hour
)

// 登录失败记录
type loginFailure struct {
	count        int
	lastFailedAt time.Time
	lockedUntil  time.Time
}

// refreshToken 会话
type refreshSession struct {
	userID    uint
	expiresAt time.Time
}

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshExpired      = errors.New("refresh token expired")
)

// recordLoginFailure 记录失败并返回是否锁定及重试时间
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

// generateRefreshToken 生成高熵的 refresh token
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// issueTokens 签发新的 access/refresh，并为用户旋转旧 refresh
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

// validateRefresh 校验 refresh token 并返回用户
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

// rotateRefresh 删除旧 token
func (s *Server) dropRefresh(token string) {
	s.refreshMu.Lock()
	delete(s.refreshStore, token)
	s.refreshMu.Unlock()
}
