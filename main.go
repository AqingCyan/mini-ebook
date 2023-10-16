package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mini-ebook/internal/repository"
	"mini-ebook/internal/repository/dao"
	"mini-ebook/internal/service"
	"mini-ebook/internal/web"
	"mini-ebook/internal/web/middleware"
	"strings"
	"time"
)

func main() {
	db := initDB()

	server := initWebServer()

	initUserHandler(db, server)

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/mini_ebook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	// cors middleware
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

	useJWT(server)

	return server
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
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
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("FrNCWQJKwK0W3yATzClayboYmU700J5B"), []byte("Qg2fflCzmy2bn5dNOsUVHtCyJKGbK4u5"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
