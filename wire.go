//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"mini-ebook/internal/repository"
	"mini-ebook/internal/repository/cache"
	"mini-ebook/internal/repository/dao"
	"mini-ebook/internal/service"
	"mini-ebook/internal/web"
	"mini-ebook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 初始化 第三方依赖
		ioc.InitRedis, ioc.InitDB,

		// 初始化 DAO 依赖
		dao.NewUserDao,

		// 初始化 cache 依赖
		cache.NewUserCache, cache.NewCodeCache,

		// 初始化 repository 依赖
		repository.NewUserRepository, repository.NewCodeRepository,

		// 初始化 service 依赖
		ioc.InitSMSService, service.NewUserService, service.NewCodeService,

		// 初始化 handler 依赖
		web.NewUserHandler,

		// 初始化 middleware web
		ioc.InitGinMiddlewares, ioc.InitWebServer,
	)
	return gin.Default()
}
