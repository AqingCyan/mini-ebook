package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now())

	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}

		session := sessions.Default(ctx)
		userId := session.Get("userId")

		println(userId)

		if userId == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 如何知道需要刷新 session
		now := time.Now()
		const updateTimeKey = "update_time"
		val := session.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)

		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			session.Set(updateTimeKey, now)
			session.Set("userId", userId)
			session.Options(sessions.Options{MaxAge: 30})
			err := session.Save()
			if err != nil {
				// 打日志
				fmt.Println(err)
			}
		}
	}
}
