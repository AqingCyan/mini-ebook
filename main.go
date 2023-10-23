package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	sessionRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mini-ebook/config"
	"mini-ebook/internal/repository"
	"mini-ebook/internal/repository/cache"
	"mini-ebook/internal/repository/dao"
	"mini-ebook/internal/service"
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
	initUserHandler(db, server, redisClient)

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

func initUserHandler(db *gorm.DB, server *gin.Engine, redisClient redis.Cmdable) {
	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(redisClient)
	ur := repository.NewUserRepository(ud, uc)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func useJWT(server *gin.Engine) {
	login := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
	// session setter middleware
	login := &middleware.LoginMiddlewareBuilder{}
	store, err := sessionRedis.NewStore(
		16, "tcp", "localhost:6379", "",
		[]byte("FrNCWQJKwK0W3yATzClayboYmU700J5B"), []byte("Qg2fflCzmy2bn5dNOsUVHtCyJKGbK4u5"),
	)
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
