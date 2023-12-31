package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"mini-ebook/internal/web"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" {
			return
		}

		// 根据约定，token 在 Authorization 头部 "Bearer TokenContent......"
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JwtKey, nil
		})
		if err != nil {
			println(err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			// token 解析成功但是非法或者过期了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			// 后期讲到了监控告警的时候，这个地方要埋点
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt
		if expireTime.Before(time.Now()) {
			// Valid 其实已经判定了过期，但出于保险，还是在这里做一个处理
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 剩余 20 分钟过期的时候刷新 Token
		if expireTime.Sub(time.Now()) < time.Minute*20 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 60))
			tokenStr, err = token.SignedString(web.JwtKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				log.Println(err)
			}
		}
		ctx.Set("user", uc)
	}
}
