package todo

import "time"

// 常量定义，避免魔法数字

// 分页相关常量
const (
	DefaultPageSize = 20  // 默认每页条数
	MaxPageSize     = 100 // 最大每页条数
	MinPageSize     = 1   // 最小每页条数
	DefaultPage     = 1   // 默认页码
)

// 限流相关常量
const (
	DefaultRateLimitWindow = 1 * time.Minute // 默认限流窗口
	DefaultRateLimitCount  = 100             // 默认限流次数
	CleanupInterval        = 1 * time.Hour   // 清理过期数据的间隔
)

// 认证相关常量
const (
	DefaultJWTTTL       = 24 * time.Hour     // 默认 JWT 过期时间
	DefaultRefreshTTL   = 7 * 24 * time.Hour // 默认刷新令牌过期时间
	LoginFailureWindow  = 15 * time.Minute   // 登录失败记录窗口
	MaxLoginFailures    = 5                  // 最大登录失败次数
	AccountLockDuration = 15 * time.Minute   // 账户锁定时长
)

// 密码策略常量
const (
	MinPasswordLength = 8  // 最小密码长度
	MaxPasswordLength = 72 // 最大密码长度（bcrypt 限制）
)

// HTTP 超时常量
const (
	DefaultReadTimeout  = 15 * time.Second
	DefaultWriteTimeout = 15 * time.Second
	DefaultIdleTimeout  = 60 * time.Second
	ShutdownTimeout     = 30 * time.Second // 优雅关闭超时时间
)
