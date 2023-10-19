package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"mini-ebook/internal/domain"
	"mini-ebook/internal/service"
	"net/http"
	"time"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
)

var JwtKey = []byte("FrNCWQJKwK0W3yATzClayboYmU700J5A")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}

type UserHandler struct {
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
	}
}

func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/users")

	group.POST("/signup", uh.Signup)
	//group.POST("/login", uh.Login)
	group.POST("/login", uh.LoginJWT)
	group.POST("/edit", uh.Edit)
	group.GET("/profile", uh.Profile)
}

func (uh *UserHandler) Signup(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignupReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := uh.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码不匹配")
		return
	}

	isPassword, err := uh.passwordRegExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	err = uh.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "注册成功")
	case errors.Is(err, service.ErrDuplicateEmail):
		ctx.String(http.StatusOK, "邮箱冲突，请尝试别的邮箱")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (uh *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	u, err := uh.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		session := sessions.Default(ctx)
		session.Set("userId", u.Id)
		session.Options(sessions.Options{MaxAge: 3600, HttpOnly: true})
		err := session.Save() // must call "save func"
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "用户不存在或是密码不正确")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (uh *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	u, err := uh.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		uc := UserClaims{
			Uid: u.Id,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)), // 60 分钟过期
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		tokenStr, err := token.SignedString(JwtKey)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
		}
		ctx.Header("x-jwt-token", tokenStr)
		ctx.String(http.StatusOK, "登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "用户不存在或是密码不正确")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (uh *UserHandler) Edit(ctx *gin.Context) {
}

func (uh *UserHandler) Profile(ctx *gin.Context) {
	//us := ctx.MustGet("user").(UserClaims)
	ctx.String(http.StatusOK, "这是 Profile")
}
