package service

import (
	"context"
	"fmt"
	"math/rand"
	"mini-ebook/internal/repository"
	"mini-ebook/internal/service/sms"
)

type CodeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

// Send 生成一个随机验证码，并发送
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 要开始发送验证码
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

// Verify 验证验证码
func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return true, nil
}

func (svc *CodeService) generate() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
