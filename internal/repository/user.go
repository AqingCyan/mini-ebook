package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"mini-ebook/internal/domain"
	"mini-ebook/internal/repository/cache"
	"mini-ebook/internal/repository/dao"
	"time"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) UpdateUserInfo(ctx *gin.Context, u domain.User) error {
	return repo.dao.UpdateByUserId(ctx, repo.toDaoEntity(u))
}

func (repo *UserRepository) FindById(ctx *gin.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	// 只要 err 为 nil 就返回，err 不为 nil 就去查询数据库
	// 而 err 有两种可能：
	// 1、key 不存在，说明 redis 是正常的。
	// 2、访问 redis 有问题，可能是网络有问题，也有可能是 redis 本身就崩溃了。
	if err == nil {
		return du, nil
	}

	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}

	du = repo.toDomain(u)
	err = repo.cache.Set(ctx, du)
	// redis 可能是网络问题，或者本身崩掉了
	if err != nil {
		log.Println(err)
	}

	return du, nil
}

/* --- 一些内部用的工具方法 --- */

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}
}

func (repo *UserRepository) toDaoEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}
