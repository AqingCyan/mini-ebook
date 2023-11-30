package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mini-ebook/internal/domain"
	"time"
)

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

// NewUserCache 创建基于 Redis 的缓存
func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

// Get 获取缓存中的 User
func (c RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.key(uid)
	// 读取缓存后反序列化
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}

	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

// Set 将 User 序列化后设置到缓存中
func (c RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := c.key(du.Id)
	// 序列化后缓存
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}

	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

// key 读写 RedisUserCache 缓存时候的键
func (c *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}
