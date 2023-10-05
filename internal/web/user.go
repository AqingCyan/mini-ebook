package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mini-ebook/internal/domain"
	"mini-ebook/internal/service"
	"net/http"
	"strings"
	"time"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
)

var validationErrors = map[string]string{
	"max": "超过最大长度限制",
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
		session.Options(sessions.Options{
			MaxAge:   3600 * 12,
			HttpOnly: true,
		})
		err := session.Save() // 必须调用 Save 方法才能保证 session 设置的字段生效
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

func (uh *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		AboutMe  string `json:"aboutMe" validate:"max=200"`
		Birthday string `json:"birthday"`
		Nickname string `json:"nickname" validate:"max=20"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	var validate *validator.Validate
	validate = validator.New()

	if err := validate.Struct(req); err != nil {
		var e []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("%s %s", err.Namespace(), validationErrors[err.ActualTag()])
			e = append(e, errorMessage)
		}
		ctx.String(http.StatusBadRequest, strings.Join(e, "; "))
		return
	}

	t, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "生日格式异常")
		return
	}

	if uId, ok := checkUserId(ctx); ok {
		err = uh.svc.UpdateUserInfo(ctx, domain.User{
			Id:       uId,
			Nickname: req.Nickname,
			Birthday: t,
			AboutMe:  req.AboutMe,
		})
		if err != nil {
			ctx.String(http.StatusOK, "更新失败")
			return
		}
	} else {
		ctx.String(http.StatusOK, "用户信息错误")
	}

	ctx.String(http.StatusOK, "更新成功")
}

func (uh *UserHandler) Profile(ctx *gin.Context) {
	if uId, ok := checkUserId(ctx); ok {
		u, err := uh.svc.FindInfoByUserId(ctx, uId)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}

		type User struct {
			Nickname string `json:"nickname"`
			Email    string `json:"email"`
			AboutMe  string `json:"aboutMe"`
			Birthday string `json:"birthday"`
		}
		ctx.JSON(http.StatusOK, User{
			Nickname: u.Nickname,
			Email:    u.Email,
			AboutMe:  u.AboutMe,
			Birthday: u.Birthday.Format(time.DateOnly),
		})
	} else {
		ctx.String(http.StatusOK, "用户信息错误")
	}
}

/* ---工具方法--- */

func checkUserId(ctx *gin.Context) (int64, bool) {
	session := sessions.Default(ctx)
	uId, ok := session.Get("userId").(int64)
	if !ok {
		return 0, false
	}
	return uId, true
}
