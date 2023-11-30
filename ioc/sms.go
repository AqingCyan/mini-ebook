package ioc

import (
	"mini-ebook/internal/service/sms"
	"mini-ebook/internal/service/sms/localSms"
)

func InitSMSService() sms.Service {
	return localSms.NewServcie()
}
