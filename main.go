package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	// cors 跨域中间件
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "aqingcyan.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// session 设置的中间件
	login := &middleware.LoginMiddlewareBuilder{}
	store := cookie.NewStore([]byte("secret")) // 需要先初始化 session
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())

	return server
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}
