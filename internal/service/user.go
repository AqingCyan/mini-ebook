package service

import (
	"context"
	"mini-ebook/internal/domain"
	"mini-ebook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	return svc.repo.Create(ctx, u)
}
