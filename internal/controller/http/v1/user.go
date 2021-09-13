package v1

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	goJwt "github.com/golang-jwt/jwt"
	dto2 "github.com/scriptscat/cloudcat/internal/domain/safe/dto"
	service2 "github.com/scriptscat/cloudcat/internal/domain/safe/service"
	service3 "github.com/scriptscat/cloudcat/internal/domain/system/service"
	"github.com/scriptscat/cloudcat/internal/domain/user/dto"
	"github.com/scriptscat/cloudcat/internal/domain/user/service"
	"github.com/scriptscat/cloudcat/internal/pkg/errs"
	"github.com/scriptscat/cloudcat/internal/pkg/httputils"
	"github.com/scriptscat/cloudcat/pkg/middleware/jwt"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
)

const JwtAuthMaxAge = 432000
const JwtAutoRenew = 259200

type User struct {
	service.User
	oauthSvc service.OAuth
	sender   service3.Sender
	safe     service2.Safe
	jwtToken string
}

func NewUser(jwtToken string, svc service.User, oauthSvc service.OAuth, safe service2.Safe, sender service3.Sender) *User {
	return &User{jwtToken: jwtToken, User: svc, safe: safe, sender: sender, oauthSvc: oauthSvc}
}

// @Summary     用户
// @Description 用户登录
// @ID          login
// @Tags  	    user
// @Produce     json
// @Accept      application/x-www-form-urlencoded
// @Param       account formData string true "邮箱/手机"
// @Param       password formData string true "登录密码"
// @Success     200
// @Failure     400 {object} errs.JsonRespondError
// @Router      /user/login [post]
func (u *User) login(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		login := &dto.Login{}
		err := ctx.ShouldBind(login)
		if err != nil {
			return err
		}
		var resp *dto.UserInfo
		if err := u.safe.Rate(&dto2.SafeUserinfo{
			Identifier: login.Account,
		}, &dto2.SafeRule{
			Name:        "user-login",
			Description: "用户登录失败",
			PeriodCnt:   5,
			Period:      300 * time.Second,
		}, func() error {
			resp, err = u.Login(login)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
		return u.oauthHandle(ctx, &dto.OAuthRespond{
			UserInfo: resp,
			IsBind:   true,
		})
	})
}

// @Summary     用户
// @Description 用户注册
// @ID          register
// @Tags  	    user
// @Produce     json
// @Accept      application/x-www-form-urlencoded
// @Param       username formData string true "用户名"
// @Param       email formData string true "邮箱"
// @Param       password formData string true "登录密码"
// @Param       repassword formData string true "再输入一次登录密码"
// @Param       email_verify_code formData string false "邮箱验证码"
// @Param       inv_code formData string false "邀请码"
// @Success     200
// @Failure     400 {object} errs.JsonRespondError
// @Router      /user/register [post]
func (u *User) register(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		register := &dto.Register{}
		err := ctx.ShouldBind(register)
		if err != nil {
			return err
		}
		var resp *dto.UserInfo
		if err := u.safe.Limit(&dto2.SafeUserinfo{
			IP: ctx.ClientIP(),
		}, &dto2.SafeRule{
			Name:        "user-register",
			Description: "用户注册失败",
			PeriodCnt:   5,
			Period:      24 * time.Hour,
		}, func() error {
			resp, err = u.User.Register(register)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
		return u.oauthHandle(ctx, &dto.OAuthRespond{
			UserInfo: resp,
			IsBind:   true,
		})
	})
}

// @Summary     用户
// @Description 请求邮箱验证码
// @ID          request-email-code
// @Tags  	    user
// @Produce     json
// @Accept      application/x-www-form-urlencoded
// @Param       email formData string true "邮箱"
// @Success     200
// @Failure     400 {object} errs.JsonRespondError
// @Router      /user/request-email-code [post]
func (u *User) requestEmailCode(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		email := ctx.PostForm("email")
		if email == "" {
			return errs.NewBadRequestError(1001, "邮箱不能为空")
		}
		return u.safe.Limit(&dto2.SafeUserinfo{
			Identifier: email,
		}, &dto2.SafeRule{
			Name:        "register-email-code",
			Description: "请求邮箱验证码失败",
			Interval:    30,
			PeriodCnt:   5,
			Period:      24 * time.Hour,
		}, func() error {
			code, err := u.RequestRegisterEmailCode(email)
			if err != nil {
				return err
			}
			return u.sender.SendEmail(email, "注册验证码", "您的注册验证码为:"+code.Code+" 请于5分钟内输入", "text/html")
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
func (u *User) bbsOAuth(ctx *gin.Context) {
	redirect := fmt.Sprintf("%s/api/v1/auth/bbs/callback?redirect=%s", ctx.Request.Header.Get("Origin"), url2.PathEscape(ctx.Query("redirect")))
	url, err := u.oauthSvc.RedirectOAuth(redirect, "bbs")
	if err != nil {
		httputils.HandleError(ctx, err)
		return
	}
	ctx.Redirect(http.StatusFound, url)
}

// @Summary     用户
// @Description 微信oauth2.0登录
// @ID          wechat-login
// @Tags  	    user
// @Success     302
// @Failure     400 {object} errs.JsonRespondError
// @Router      /auth/wechat [post]
func (u *User) wechatOAuth(ctx *gin.Context) {
	redirect := fmt.Sprintf("%s/api/v1/auth/wechat/callback?redirect=%s", ctx.Request.Header.Get("Origin"), url2.PathEscape(ctx.Query("redirect")))
	url, err := u.oauthSvc.RedirectOAuth(redirect, "wechat")
	if err != nil {
		httputils.HandleError(ctx, err)
		return
	}
	ctx.Redirect(http.StatusFound, url)
}

func (u *User) bbsOAuthCallback(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		code := ctx.Query("code")
		if code == "" {
			return errs.NewBadRequestError(1001, "code不能为空")
		}
		resp, err := u.oauthSvc.BBSOAuthLogin(code)
		if err != nil {
			return err
		}
		return u.oauthHandle(ctx, resp)
	})
}

func (u *User) wechatOAuthCallback(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		code := ctx.Query("code")
		if code == "" {
			return errs.NewBadRequestError(1001, "code不能为空")
		}
		resp, err := u.oauthSvc.WechatAuthLogin(code)
		if err != nil {
			return err
		}
		return u.oauthHandle(ctx, resp)
	})
}

func (u *User) wechatHandle(ctx *gin.Context) {
	wc, err := u.oauthSvc.GetWechat()
	if err != nil {
		httputils.HandleError(ctx, err)
		return
	}
	if wc.Officialaccount == nil {
		return
	}
	s := wc.Officialaccount.GetServer(ctx.Request, ctx.Writer)
	s.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		switch msg.MsgType {
		case message.MsgTypeEvent:
			param := ""
			switch msg.Event {
			case message.EventSubscribe:
				param = strings.TrimSuffix(msg.EventKey, "qrscene_")
				if param == msg.EventKey {
					return nil
				}
			case message.EventScan:
				param = msg.EventKey
			}
			code := strings.Split(param, "_")
			if len(code) != 2 {
				return nil
			}
			if code[0] == "login" {
				if err := u.oauthSvc.WechatScanLogin(string(msg.FromUserName), code[1]); err != nil {
					logrus.Errorf("wx login message handler: %v", err)
					return nil
				}
			}
		}

		return nil
	})

	if err := s.Serve(); err != nil {
		logrus.Errorf("wx message handler: %v", err)
		return
	}

	if wc.ReverseProxy != "" {
		url, _ := url2.Parse(wc.ReverseProxy)
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = url.Scheme
				req.URL.Host = url.Host
				req.URL.Path, req.URL.RawPath = url.Path, url.RawPath
				if url.RawQuery == "" || req.URL.RawQuery == "" {
					req.URL.RawQuery = url.RawQuery + req.URL.RawQuery
				} else {
					req.URL.RawQuery = url.RawQuery + "&" + req.URL.RawQuery
				}
				if _, ok := req.Header["User-Agent"]; !ok {
					// explicitly disable User-Agent so it's not set to default value
					req.Header.Set("User-Agent", "")
				}
			},
		}
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
		return
	}

	if err := s.Send(); err != nil {
		logrus.Errorf("wx message send: %v", err)
	}
}

func (u *User) oauthHandle(ctx *gin.Context, resp *dto.OAuthRespond) interface{} {
	if !resp.IsBind {
		// 跳转到注册页面
		return errs.NewBadRequestError(1002, "账号未注册,请先注册后绑定三方平台")
	}
	tokenString, err := jwt.GenJwt([]byte(u.jwtToken), goJwt.MapClaims{
		"uid":      resp.UserInfo.ID,
		"username": resp.UserInfo.Nickname,
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

func (u *User) Register(r *gin.RouterGroup) {
	rg := r.Group("/user")
	rg.POST("/login", u.login)
	rg.POST("/register", u.register)
	rg.POST("/request-email-code", u.requestEmailCode)

	rg = r.Group("/auth")
	rg.POST("/bbs", u.bbsOAuth)
	rg.GET("/bbs/callback", u.bbsOAuthCallback)
	rg.POST("/wechat", u.wechatOAuth)
	//rg.GET("/wechat/callback", s.wechatOAuthCallback)
	rg.Any("/wechat/handle", u.wechatHandle)

}
