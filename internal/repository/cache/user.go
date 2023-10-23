package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"mini-ebook/internal/domain"
	"time"
)

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

// Get 获取缓存中的 User
func (c UserCache) Get(ctx *gin.Context, uid int64) (domain.User, error) {
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
func (c UserCache) Set(ctx *gin.Context, du domain.User) error {
	key := c.key(du.Id)
	// 序列化后缓存
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}

	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

// key 读写 UserCache 缓存时候的键
func (c *UserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
