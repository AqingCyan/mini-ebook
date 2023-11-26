package repository

import (
	"context"
	"mini-ebook/internal/repository/cache"
)

type CodeRepository struct {
	cache cache.CodeCache
}

func (c *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}
