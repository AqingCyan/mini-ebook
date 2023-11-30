package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"mini-ebook/internal/domain"
	"mini-ebook/internal/repository"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户不存在或是密码不正确")
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	UpdateUserInfo(ctx context.Context, user domain.User) error
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	FindInfoByUserId(ctx context.Context, uid int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	// 检查密码是否匹配
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) UpdateUserInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateUserInfo(ctx, user)
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *userService) FindInfoByUserId(ctx context.Context, uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 先找一下，我们认为大部分用户是已经存在的用户
	_, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		// 有两种情况
		// err == nil, u 是可用的
		// err != nil, 系统错误
		return domain.User{}, err
	}

	// 如果没找到该用户
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 两种可能：一种是 err 恰好是唯一索引冲突（phone）；一种是 err 不为 nil，系统错误
	if err != nil && !errors.Is(err, repository.ErrDuplicateUser) {
		return domain.User{}, err
	}

	// 要么 err == nil，要么 ErrDuplicateUser 也代表用户存在
	// 考虑到生产环境有主从延迟，刚插入的数据不一定能立马查出来，所以这里理论上讲应该强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}
