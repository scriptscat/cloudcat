package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	goJwt "github.com/golang-jwt/jwt"
	"github.com/scriptscat/cloudcat/internal/controller/http/v1/dto/request"
	"github.com/scriptscat/cloudcat/internal/domain/user/service"
	"github.com/scriptscat/cloudcat/internal/pkg/errs"
	"github.com/scriptscat/cloudcat/internal/pkg/httputils"
	"github.com/scriptscat/cloudcat/pkg/middleware/jwt"
)

const JwtAuthMaxAge = 432000
const JwtAutoRenew = 259200

type User struct {
	jwtToken string
	svc      service.User
}

func NewUser(jwtToken string, svc service.User) *User {
	return &User{jwtToken: jwtToken, svc: svc}
}

// @Summary     用户
// @Description 用户登录
// @ID          user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Accept      application/x-www-form-urlencoded
// @Param       name formData string false "用户名"
// @Param       email formData string false "邮箱"
// @Param       password formData string true "登录密码"
// @Success     200 {object} repository.ScriptCatInfo
// @Failure     400 {object} errs.JsonRespondError
// @Router      /user/login [post]
func (s *User) login(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		script := &request.UserLogin{}
		if err := ctx.ShouldBind(script); err != nil {
			return err
		}
		return nil
	})
}

// @Summary     用户
// @Description 用户注册
// @ID          user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Accept      application/x-www-form-urlencoded
// @Param       name formData string true "用户名"
// @Param       email formData string true "邮箱"
// @Param       password formData string true "登录密码"
// @Param       rePassword formData string true "再输入一次登录密码"
// @Success     200 {object} repository.ScriptCatInfo
// @Failure     400 {object} errs.JsonRespondError
// @Router      /user/login [post]
func (s *User) register(ctx *gin.Context) {

}

// @Summary     用户
// @Description 论坛oauth2.0登录
// @ID          user
// @Tags  	    user
// @Success     302
// @Failure     400 {object} errs.JsonRespondError
// @Router      /auth/bbs [post]
func (s *User) bbsOAuth(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		url, err := s.svc.RedirectOAuth(ctx.Request.URL.String(), "bbs")
		if err != nil {
			return err
		}
		ctx.Redirect(http.StatusFound, url)
		return nil
	})
}

// @Summary     用户
// @Description 微信oauth2.0登录
// @ID          user
// @Tags  	    user
// @Success     302
// @Failure     400 {object} errs.JsonRespondError
// @Router      /auth/wechat [post]
func (s *User) wechatOAuth(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		url, err := s.svc.RedirectOAuth(ctx.Request.URL.String(), "wechat")
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
		if userInfo, err := s.svc.BBSOAuthLogin(code); err != nil {
			return err
		} else {
			tokenString, err := jwt.GenJwt([]byte(s.jwtToken), goJwt.MapClaims{
				"uid":      userInfo.ID,
				"username": userInfo.Username,
				"email":    userInfo.Email,
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
	})
}

func (s *User) wechatOAuthCallback(ctx *gin.Context) {

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
