package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/users")

	group.POST("/signup", uh.Signup)
	group.POST("/login", uh.Login)
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
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusBadRequest, "非法邮箱格式")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusBadRequest, "两次密码不匹配")
		return
	}

	isPassword, err := uh.passwordRegExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusBadRequest, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	ctx.String(http.StatusOK, "Hello World")
}

func (uh *UserHandler) Login(ctx *gin.Context) {

}

func (uh *UserHandler) Edit(ctx *gin.Context) {

}

func (uh *UserHandler) Profile(ctx *gin.Context) {

}