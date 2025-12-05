// Package cache 提供 Redis 缓存功能
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// RedisCache Redis 缓存客户端
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建 Redis 缓存
func NewRedisCache(cfg RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     50,
		MinIdleConns: 10,
		MaxRetries:   3,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis 连接失败: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Close 关闭连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// Set 设置缓存
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get 获取缓存
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) bool {
	return c.client.Exists(ctx, key).Val() > 0
}

// SetNX 仅当键不存在时设置 (分布式锁)
func (c *RedisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return c.client.SetNX(ctx, key, data, ttl).Result()
}

// 限流脚本 (固定窗口)
const rateLimitScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCR", key)
if current == 1 then
    redis.call("PEXPIRE", key, window)
end
if current > limit then
    return 0
end
return 1
`

// RateLimit 限流检查
// limit: 窗口内最大请求数
// windowMs: 窗口时间(毫秒)
func (c *RedisCache) RateLimit(ctx context.Context, key string, limit int, windowMs int) (bool, error) {
	result, err := c.client.Eval(ctx, rateLimitScript, []string{key}, limit, windowMs).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// Client 返回底层客户端 (高级用法)
func (c *RedisCache) Client() *redis.Client {
	return c.client
}
