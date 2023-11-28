package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mini-ebook/config"
	"mini-ebook/internal/repository"
	"mini-ebook/internal/repository/cache"
	"mini-ebook/internal/repository/dao"
	"mini-ebook/internal/service"
	"mini-ebook/internal/service/sms"
	"mini-ebook/internal/service/sms/localSms"
	"mini-ebook/internal/web"
	"mini-ebook/internal/web/middleware"
	"mini-ebook/pkg/ginx/middleware/ratelimit"
	"strings"
	"time"
)

func main() {
	db := initDB()
	redisClient := redis.NewClient(&redis.Options{Addr: config.Config.Redis.Addr})

	server := initWebServer(redisClient)
	userDB := dao.NewUserDao(db)
	userSvc := initUserSvc(userDB, redisClient)
	codeSvc := initCodeSvc(redisClient)
	initUserHandler(userSvc, codeSvc, server)

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initWebServer(redisClient redis.Cmdable) *gin.Engine {
	server := gin.Default()

	// cors 跨域中间件
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"x-jwt-token"}, // 允许前端访问服务响应的头部
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "aqingcyan.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// 如果为了压测，得去掉下面的限流
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 10).Build())

	useJWT(server)

	return server
}

func initUserHandler(userSvc *service.UserService, codeSvc *service.CodeService, server *gin.Engine) {
	hdl := web.NewUserHandler(userSvc, codeSvc)
	hdl.RegisterRoutes(server)
}

func initUserSvc(ud *dao.UserDao, redisClient redis.Cmdable) *service.UserService {
	uc := cache.NewUserCache(redisClient)
	ur := repository.NewUserRepository(ud, uc)
	return service.NewUserService(ur)
}

func initCodeSvc(redisClient redis.Cmdable) *service.CodeService {
	cc := cache.NewCodeCache(redisClient)
	crepo := repository.NewCodeRepository(cc)
	return service.NewCodeService(crepo, initMemorySms())
}

func initMemorySms() sms.Service {
	return localSms.NewServcie()
}

func useJWT(server *gin.Engine) {
	login := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}
