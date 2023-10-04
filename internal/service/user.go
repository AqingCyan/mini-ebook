package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"mini-ebook/internal/domain"
	"mini-ebook/internal/repository"
)

var ErrDuplicateEmail = repository.ErrDuplicateEmail

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return svc.repo.Create(ctx, u)
}
