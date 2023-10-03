package main

import (
	"github.com/gin-gonic/gin"
	"mini-ebook/internal/web"
)

func main() {
	server := gin.Default()

	hdl := web.NewUserHandler()
	hdl.RegisterRoutes(server)

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
