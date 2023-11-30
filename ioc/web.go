package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"mini-ebook/internal/web"
	"mini-ebook/internal/web/middleware"
	"mini-ebook/pkg/ginx/middleware/ratelimit"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// cors 跨域中间件
		cors.New(cors.Config{
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
		}),

		// 限流（如果为了压测，得去掉限流）
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),

		// JWT 验权
		(&middleware.LoginJWTMiddlewareBuilder{}).CheckLogin(),
	}
}
