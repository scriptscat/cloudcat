package v1

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	dto2 "github.com/scriptscat/cloudcat/internal/domain/safe/dto"
	service2 "github.com/scriptscat/cloudcat/internal/domain/safe/service"
	service3 "github.com/scriptscat/cloudcat/internal/domain/system/service"
	"github.com/scriptscat/cloudcat/internal/domain/user/dto"
	"github.com/scriptscat/cloudcat/internal/domain/user/service"
	"github.com/scriptscat/cloudcat/internal/pkg/errs"
	"github.com/scriptscat/cloudcat/internal/pkg/httputils"
	"github.com/scriptscat/cloudcat/pkg/cache"
	"github.com/scriptscat/cloudcat/pkg/middleware/token"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
)

type Auth struct {
	service.User
	oauthSvc service.OAuth
	sender   service3.Sender
	safe     service2.Safe
	cache    cache.Cache
}

func NewAuth(cache cache.Cache, svc service.User, oauthSvc service.OAuth, safe service2.Safe, sender service3.Sender) *Auth {
	return &Auth{cache: cache, User: svc, safe: safe, sender: sender, oauthSvc: oauthSvc}
}

// @Summary     用户
// @Description 用户登录
// @ID          login
// @Tags  	    user
// @Produce     json
// @Accept      x-www-form-urlencoded
// @Param       username formData string true "用户名/邮箱"
// @Param       password formData string true "登录密码"
// @Param       auto_login formData bool false "自动登录"
// @Success     200
// @Failure     400 {object} errs.JsonRespondError
// @Router      /account/login [post]
func (a *Auth) login(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		login := &dto.Login{}
		err := ctx.ShouldBind(login)
		if err != nil {
			return err
		}
		var resp *dto.UserInfo
		if err := a.safe.Rate(&dto2.SafeUserinfo{
			IP:         ctx.ClientIP(),
			Identifier: login.Username,
		}, &dto2.SafeRule{
			Name:        "user-login",
			Description: "用户登录失败",
			PeriodCnt:   5,
			Period:      300 * time.Second,
		}, func() error {
			resp, err = a.Login(login)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
		return a.oauthHandle(ctx, &dto.OAuthRespond{
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
// @Accept      x-www-form-urlencoded
// @Param       username formData string true "用户名"
// @Param       email formData string true "邮箱"
// @Param       password formData string true "登录密码"
// @Param       repassword formData string true "再输入一次登录密码"
// @Param       email_verify_code formData string false "邮箱验证码"
// @Param       inv_code formData string false "邀请码"
// @Success     200
// @Failure     400 {object} errs.JsonRespondError
// @Router      /account/register [post]
func (a *Auth) register(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		register := &dto.Register{}
		err := ctx.ShouldBind(register)
		if err != nil {
			return err
		}
		var resp *dto.UserInfo
		if err := a.safe.Limit(&dto2.SafeUserinfo{
			IP: ctx.ClientIP(),
		}, &dto2.SafeRule{
			Name:        "user-register",
			Description: "用户注册失败",
			PeriodCnt:   5,
			Period:      24 * time.Hour,
		}, func() error {
			resp, err = a.User.Register(register)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
		return a.oauthHandle(ctx, &dto.OAuthRespond{
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
// @Accept      x-www-form-urlencoded
// @Param       email formData string true "邮箱"
// @Success     200
// @Failure     400 {object} errs.JsonRespondError
// @Router      /account/register/request-email-code [post]
func (a *Auth) requestEmailCode(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		email := ctx.PostForm("email")
		if email == "" {
			return errs.NewBadRequestError(1001, "邮箱不能为空")
		}
		return a.safe.Limit(&dto2.SafeUserinfo{
			Identifier: email,
		}, &dto2.SafeRule{
			Name:        "register-email-code",
			Description: "请求邮箱验证码失败",
			Interval:    30,
			PeriodCnt:   5,
			Period:      24 * time.Hour,
		}, func() error {
			code, err := a.RequestRegisterEmailCode(email)
			if err != nil {
				return err
			}
			return a.sender.SendEmail(email, "注册验证码", "您的注册验证码为:"+code.Code+" 请于5分钟内输入", "text/html")
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
func (a *Auth) bbsOAuth(ctx *gin.Context) {
	redirect := fmt.Sprintf("%s/api/v1/auth/bbs/callback?redirect=%s", ctx.Request.Header.Get("Origin"), url.PathEscape(ctx.Query("redirect")))
	url, err := a.oauthSvc.RedirectOAuth(redirect, "bbs")
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
func (a *Auth) wechatOAuth(ctx *gin.Context) {
	redirect := fmt.Sprintf("%s/api/v1/auth/wechat/callback?redirect=%s", ctx.Request.Header.Get("Origin"), url.PathEscape(ctx.Query("redirect")))
	url, err := a.oauthSvc.RedirectOAuth(redirect, "wechat")
	if err != nil {
		httputils.HandleError(ctx, err)
		return
	}
	ctx.Redirect(http.StatusFound, url)
}

// @Summary     用户
// @Description 请求微信登录二维码
// @ID          wechat-request
// @Tags  	    user
// @Success     200 {object} dto.WechatScanLogin
// @Failure     400 {object} errs.JsonRespondError
// @Router      /auth/wechat/request [post]
func (a *Auth) wechatRequest(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		ret, err := a.oauthSvc.WechatScanLoginRequest()
		if err != nil {
			return err
		}
		return ret
	})
}

// @Summary     用户
// @Description 查询微信扫码状态
// @ID          wechat-status
// @Tags  	    user
// @Param       redirect_uri query string false "重定向链接"
// @param       code formData string true "查询code"
// @Success     200 {string} json "token"
// @Success     302
// @Failure     400 {object} errs.JsonRespondError
// @Failure     404 {object} errs.JsonRespondError
// @Router      /auth/wechat/status [post]
func (a *Auth) wechatStatus(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		code := ctx.PostForm("code")
		if code == "" {
			return errs.NewBadRequestError(1000, "查询状态码不能为空")
		}
		ret, err := a.oauthSvc.WechatScanLoginStatus(code)
		if err != nil {
			return err
		}
		return a.oauthHandle(ctx, ret)
	})
}

func (a *Auth) bbsOAuthCallback(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		code := ctx.Query("code")
		if code == "" {
			return errs.NewBadRequestError(1001, "code不能为空")
		}
		resp, err := a.oauthSvc.BBSOAuthLogin(code)
		if err != nil {
			return err
		}
		return a.oauthHandle(ctx, resp)
	})
}

func (a *Auth) wechatOAuthCallback(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		code := ctx.Query("code")
		if code == "" {
			return errs.NewBadRequestError(1001, "code不能为空")
		}
		resp, err := a.oauthSvc.WechatAuthLogin(code)
		if err != nil {
			return err
		}
		return a.oauthHandle(ctx, resp)
	})
}

func (a *Auth) wechatHandle(ctx *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logrus.Errorf("wx handel read body: %v", err)
		return
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	wc, err := a.oauthSvc.GetWechat()
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
				if err := a.oauthSvc.WechatScanLogin(string(msg.FromUserName), code[1]); err != nil {
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

	if ctx.Query("echostr") == "" && wc.ReverseProxy != "" {
		resp, err := http.Post(wc.ReverseProxy+"?"+ctx.Request.URL.RawQuery, "application/xml", bytes.NewBuffer(bodyBytes))
		ctx.Writer.WriteHeader(resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("wx write body: %v", err)
			return
		}
		ctx.Writer.Write(b)
		return
	}

	if err := s.Send(); err != nil {
		logrus.Errorf("wx message send: %v", err)
	}
}

func (a *Auth) oauthHandle(ctx *gin.Context, resp *dto.OAuthRespond) interface{} {
	if !resp.IsBind {
		// 跳转到注册页面
		return errs.NewBadRequestError(1002, "账号未注册,请先注册后绑定三方平台")
	}
	tokenString, err := token.GenToken(a.cache, gin.H{
		"uid":      strconv.FormatInt(resp.UserInfo.ID, 10),
		"username": resp.UserInfo.Username,
	})
	if err != nil {
		return err
	}
	ctx.SetCookie("token", tokenString, TokenAuthMaxAge, "/", "", false, true)
	if uri := ctx.Query("redirect_uri"); uri != "" {
		ctx.Redirect(http.StatusFound, uri)
		return nil
	}
	return gin.H{
		"uid":   resp.UserInfo.ID,
		"token": tokenString,
	}
}

func (a *Auth) Register(r *gin.RouterGroup) {
	rg := r.Group("/account")
	rg.POST("/login", a.login)
	rg.POST("/register", a.register)
	rg.POST("/register/request-email-code", a.requestEmailCode)

	rg = r.Group("/auth")
	rg.POST("/bbs", a.bbsOAuth)
	rg.GET("/bbs/callback", a.bbsOAuthCallback)
	rg.POST("/wechat", a.wechatOAuth)
	//rg.GET("/wechat/callback", s.wechatOAuthCallback)
	rg.POST("/wechat/request", a.wechatRequest)
	rg.POST("/wechat/status", a.wechatStatus)
	rg.Any("/wechat/handle", a.wechatHandle)
}
