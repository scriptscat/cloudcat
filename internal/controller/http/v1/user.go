package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goJwt "github.com/golang-jwt/jwt"
	dto2 "github.com/scriptscat/cloudcat/internal/domain/safe/dto"
	service2 "github.com/scriptscat/cloudcat/internal/domain/safe/service"
	"github.com/scriptscat/cloudcat/internal/domain/user/dto"
	"github.com/scriptscat/cloudcat/internal/domain/user/service"
	"github.com/scriptscat/cloudcat/internal/pkg/errs"
	"github.com/scriptscat/cloudcat/internal/pkg/httputils"
	"github.com/scriptscat/cloudcat/pkg/middleware/jwt"
)

const JwtAuthMaxAge = 432000
const JwtAutoRenew = 259200

type User struct {
	service.User
	safe     service2.Safe
	jwtToken string
}

func NewUser(jwtToken string, svc service.User) *User {
	return &User{jwtToken: jwtToken, User: svc}
}

// @Summary     用户
// @Description 用户登录
// @ID          login
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Accept      application/x-www-form-urlencoded
// @Param       username formData string false "用户名"
// @Param       email formData string false "邮箱"
// @Param       password formData string true "登录密码"
// @Success     200 {object}
// @Failure     400 {object} errs.JsonRespondError
// @Router      /user/login [post]
func (s *User) login(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		login := &dto.Login{}
		err := ctx.ShouldBind(login)
		if err != nil {
			return err
		}
		var resp *dto.UserInfo
		if err := s.safe.Rate(&dto2.SafeUserinfo{
			Identifier: login.Username + login.Email,
		}, &dto2.SafeRule{
			Name:        "user-login",
			Description: "用户登录失败",
			PeriodCnt:   5,
			Period:      300 * time.Second,
		}, func() error {
			resp, err = s.Login(login)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
		return s.oauthHandle(ctx, &dto.OAuthRespond{
			UserInfo: resp,
			IsBind:   true,
		})
	})
}

// @Summary     用户
// @Description 用户注册
// @ID          register
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Accept      application/x-www-form-urlencoded
// @Param       username formData string true "用户名"
// @Param       email formData string true "邮箱"
// @Param       password formData string true "登录密码"
// @Param       repassword formData string true "再输入一次登录密码"
// @Param       email_verify_code formData string false "邮箱验证码"
// @Param       inv_code formData string false "邀请码"
// @Success     200 {object}
// @Failure     400 {object} errs.JsonRespondError
// @Router      /user/login [post]
func (s *User) register(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		register := &dto.Register{}
		err := ctx.ShouldBind(register)
		if err != nil {
			return err
		}
		var resp *dto.UserInfo
		if err := s.safe.Limit(&dto2.SafeUserinfo{
			IP: ctx.ClientIP(),
		}, &dto2.SafeRule{
			Name:        "user-register",
			Description: "用户注册失败",
			PeriodCnt:   5,
			Period:      24 * time.Hour,
		}, func() error {
			resp, err = s.User.Register(register)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
		return s.oauthHandle(ctx, &dto.OAuthRespond{
			UserInfo: resp,
			IsBind:   true,
		})
	})
}

// @Summary     用户
// @Description 请求邮箱验证码
// @ID          request-email-code
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Accept      application/x-www-form-urlencoded
// @Param       email formData string true "邮箱"
// @Success     200 {object}
// @Failure     400 {object} errs.JsonRespondError
// @Router      /user/login [post]
func (s *User) requestEmailCode(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		email := ctx.PostForm("email")
		if email == "" {
			return errs.NewBadRequestError(1001, "邮箱不能为空")
		}
		register := &dto.Register{}
		err := ctx.ShouldBind(register)
		if err != nil {
			return err
		}
		return s.safe.Limit(&dto2.SafeUserinfo{
			Identifier: email,
		}, &dto2.SafeRule{
			Name:        "register-email-code",
			Description: "请求邮箱验证码失败",
			Interval:    30,
			PeriodCnt:   5,
			Period:      24 * time.Hour,
		}, func() error {
			return s.RequestRegisterEmailCode(email)
		})
	})
}

// @Summary     用户
// @Description 论坛oauth2.0登录
// @ID          bbs-login
// @Tags  	    user
// @Success     302
// @Failure     400 {object} errs.JsonRespondError
// @Router      /auth/bbs [post]
func (s *User) bbsOAuth(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		url, err := s.RedirectOAuth(ctx.Request.URL.String(), "bbs")
		if err != nil {
			return err
		}
		ctx.Redirect(http.StatusFound, url)
		return nil
	})
}

// @Summary     用户
// @Description 微信oauth2.0登录
// @ID          wechat-login
// @Tags  	    user
// @Success     302
// @Failure     400 {object} errs.JsonRespondError
// @Router      /auth/wechat [post]
func (s *User) wechatOAuth(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		url, err := s.RedirectOAuth(ctx.Request.URL.String(), "wechat")
		if err != nil {
			return err
		}
		ctx.Redirect(http.StatusFound, url)
		return nil
	})
}

func (s *User) bbsOAuthCallback(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		code := ctx.Query("code")
		if code == "" {
			return errs.NewBadRequestError(1001, "code不能为空")
		}
		resp, err := s.BBSOAuthLogin(code)
		if err != nil {
			return err
		}
		return s.oauthHandle(ctx, resp)
	})
}

func (s *User) wechatOAuthCallback(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		code := ctx.Query("code")
		if code == "" {
			return errs.NewBadRequestError(1001, "code不能为空")
		}
		resp, err := s.WechatAuthLogin(code)
		if err != nil {
			return err
		}
		return s.oauthHandle(ctx, resp)
	})
}

func (s *User) oauthHandle(ctx *gin.Context, resp *dto.OAuthRespond) interface{} {
	if !resp.IsBind {
		// 跳转到注册页面
		return errs.NewBadRequestError(1002, "账号未注册,请先注册后绑定三方平台")
	}
	tokenString, err := jwt.GenJwt([]byte(s.jwtToken), goJwt.MapClaims{
		"uid":      resp.UserInfo.ID,
		"username": resp.UserInfo.Username,
	})
	if err != nil {
		return err
	}
	ctx.SetCookie("auth", tokenString, JwtAuthMaxAge, "/", "", false, true)
	if uri := ctx.Query("redirect_uri"); uri != "" {
		ctx.Redirect(http.StatusFound, uri)
		return nil
	}
	return gin.H{
		"token": tokenString,
	}
}

func (s *User) Register(r *gin.RouterGroup) {
	rg := r.Group("/user")
	rg.POST("/login", s.login)
	rg.POST("/register", s.register)

	rg = r.Group("/auth")
	rg.POST("/bbs", s.bbsOAuth)
	rg.GET("/bbs/callback", s.bbsOAuthCallback)
	rg.POST("/wechat", s.wechatOAuth)
	rg.GET("/wechat/callback", s.wechatOAuthCallback)

}
