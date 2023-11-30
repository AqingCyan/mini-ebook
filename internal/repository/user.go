package repository

import (
	"context"
	"database/sql"
	"log"
	"mini-ebook/internal/domain"
	"mini-ebook/internal/repository/cache"
	"mini-ebook/internal/repository/dao"
	"time"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUserInfo(ctx context.Context, u domain.User) error
	UpdateNonZeroFields(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toDaoEntity(u))
}

// FindByEmail 根据邮箱查询用户信息
func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

// UpdateUserInfo 更新用户信息
func (repo *CachedUserRepository) UpdateUserInfo(ctx context.Context, u domain.User) error {
	return repo.dao.UpdateByUserId(ctx, repo.toDaoEntity(u))
}

// FindById 根据 UserId 查询用户信息
func (repo *CachedUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
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

func (repo *CachedUserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateByUserId(ctx, repo.toDaoEntity(user))
}

// FindByPhone 根据电话号码查询用户信息
func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

/* --- 一些内部用的工具方法 --- */

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}
}

func (repo *CachedUserRepository) toDaoEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}
