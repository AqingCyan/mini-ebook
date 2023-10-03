package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"mini-ebook/internal/web"
	"strings"
	"time"
)

func main() {
	server := gin.Default()
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

	hdl := web.NewUserHandler()
	hdl.RegisterRoutes(server)

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
