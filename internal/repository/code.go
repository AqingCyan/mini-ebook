package repository

import (
	"context"
	"mini-ebook/internal/repository/cache"
)

var ErrCodeVerifyTooMany = cache.ErrCodeVerifyToMany
var ErrCodeSendTooMany = cache.ErrCodeSendToMany

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (c *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
