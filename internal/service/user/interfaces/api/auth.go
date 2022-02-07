package api

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
	"github.com/scriptscat/cloudcat/internal/infrastructure/middleware/token"
	service3 "github.com/scriptscat/cloudcat/internal/infrastructure/sender"
	service2 "github.com/scriptscat/cloudcat/internal/service/safe/application"
	dto2 "github.com/scriptscat/cloudcat/internal/service/safe/domain/dto"
	application2 "github.com/scriptscat/cloudcat/internal/service/user/application"
	errs2 "github.com/scriptscat/cloudcat/internal/service/user/domain/errs"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/vo"
	"github.com/scriptscat/cloudcat/internal/service/user/interfaces/api/request"
	"github.com/scriptscat/cloudcat/pkg/cache"
	"github.com/scriptscat/cloudcat/pkg/errs"
	"github.com/scriptscat/cloudcat/pkg/httputils"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
)

type Auth struct {
	application2.User
	oauthSvc application2.OAuth
	sender   service3.Sender
	safe     service2.Safe
	cache    cache.Cache
}

func NewAuth(cache cache.Cache, svc application2.User, oauthSvc application2.OAuth, safe service2.Safe, sender service3.Sender) *Auth {
	return &Auth{cache: cache, User: svc, safe: safe, sender: sender, oauthSvc: oauthSvc}
}

// @Summary      用户
// @Description  用户登录
// @ID           login
// @Tags         user
// @Produce      json
// @Accept       x-www-form-urlencoded
// @Param        username    formData  string  true   "用户名/邮箱"
// @Param        password    formData  string  true   "登录密码"
// @Param        auto_login  formData  bool    false  "自动登录"
// @Success      200
// @Failure      400   {object}  errs.JsonRespondError
// @Router       /account/login [post]
func (a *Auth) login(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		login := &vo.Login{}
		err := ctx.ShouldBind(login)
		if err != nil {
			return err
		}
		var resp *vo.UserInfo
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
		return a.oauthHandle(ctx, &vo.OAuthRespond{
			UserInfo: resp,
			IsBind:   true,
		})
	})
}

// @Summary      用户
// @Description  登出
// @ID           logout
// @Tags         user
// @Produce      json
// @Accept       x-www-form-urlencoded
// @Success      200
// @Failure      400  {object}  errs.JsonRespondError
// @Router       /account/logout [post]
func (a *Auth) logout(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		t, ok := token.Authtoken(ctx)
		if !ok {
			return nil
		}
		ctx.SetCookie("token", "", -1, "/", "", false, true)
		token.DelToken(a.cache, t.Token)
		return "登出成功"
	})
}

// @Summary      用户
// @Description  用户注册
// @ID           register
// @Tags         user
// @Produce      json
// @Accept       x-www-form-urlencoded
// @Param        username           formData  string  true   "用户名"
// @Param        email              formData  string  true   "邮箱"
// @Param        password           formData  string  true   "登录密码"
// @Param        repassword         formData  string  true   "再输入一次登录密码"
// @Param        email_verify_code  formData  string  false  "邮箱验证码"
// @Param        inv_code           formData  string  false  "邀请码"
// @Success      200
// @Failure      400  {object}  errs.JsonRespondError
// @Router       /account/register [post]
func (a *Auth) register(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		register := &vo.Register{}
		err := ctx.ShouldBind(register)
		if err != nil {
			return err
		}
		var resp *vo.UserInfo
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
		return a.oauthHandle(ctx, &vo.OAuthRespond{
			UserInfo: resp,
			IsBind:   true,
		})
	})
}

// @Summary      用户
// @Description  请求邮箱验证码
// @ID           request-email-code
// @Tags         user
// @Produce      json
// @Accept       x-www-form-urlencoded
// @Param        email  formData  string  true  "邮箱"
// @Success      200
// @Failure      400  {object}  errs.JsonRespondError
// @Router       /account/register/request-email-code [post]
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
			code, err := a.RequestEmailCode(email, "register")
			if err != nil {
				return err
			}
			return a.sender.SendEmail(email, "注册验证码", "您的注册验证码为:"+code.Code+" 请于5分钟内输入", "text/html")
		})
	})
}

// @Summary      用户
// @Description  论坛oauth2.0登录
// @ID           bbs-login
// @Tags         user
// @Success      302
// @Failure      400  {object}  errs.JsonRespondError
// @Router       /auth/bbs [get]
func (a *Auth) bbsOAuth(ctx *gin.Context) {
	redirect := fmt.Sprintf("/api/v1/auth/bbs/callback?redirect=%s", url.PathEscape(ctx.Query("redirect")))
	url, err := a.oauthSvc.RedirectOAuth(redirect, "bbs")
	if err != nil {
		httputils.HandleError(ctx, err)
		return
	}
	ctx.Redirect(http.StatusFound, url)
}

// @Summary      用户
// @Description  微信oauth2.0登录
// @ID           wechat-login
// @Tags         user
// @Success      302
// @Failure      400  {object}  errs.JsonRespondError
// @Router       /auth/wechat [post]
func (a *Auth) wechatOAuth(ctx *gin.Context) {
	redirect := fmt.Sprintf("%s/api/v1/auth/wechat/callback?redirect=%s", ctx.Request.Header.Get("Origin"), url.PathEscape(ctx.Query("redirect")))
	url, err := a.oauthSvc.RedirectOAuth(redirect, "wechat")
	if err != nil {
		httputils.HandleError(ctx, err)
		return
	}
	ctx.Redirect(http.StatusFound, url)
}

// @Summary      用户
// @Description  请求微信登录二维码
// @ID           wechat-request
// @Tags         user
// @Success      200  {object}  vo.WechatScan
// @Failure      400  {object}  errs.JsonRespondError
// @Router       /auth/wechat/request [post]
func (a *Auth) wechatRequest(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		ret, err := a.oauthSvc.WechatScanRequest("login")
		if err != nil {
			return err
		}
		return ret
	})
}

// @Summary      用户
// @Description  绑定微信
// @ID           wechat-bind-request
// @Tags         user
// @Success      200  {object}  vo.WechatScan
// @Failure      400  {object}  errs.JsonRespondError
// @Router       /auth/bind/wechat/request [post]
func (a *Auth) wechatBindRequest(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, _ := token.UserId(ctx)
		ret, err := a.oauthSvc.WechatScanRequest("bind")
		if err != nil {
			return err
		}
		if err := a.oauthSvc.WechatScanBindCode(uid, ret.Code); err != nil {
			return err
		}
		return ret
	})
}

// @Summary      用户
// @Description  查询微信扫码状态
// @ID           wechat-status
// @Tags         user
// @param        code          formData  string  true   "查询code"
// @Success      200           {string}  json    "token"
// @Success      302
// @Failure      400  {object}  errs.JsonRespondError
// @Failure      404   {object}  errs.JsonRespondError
// @Router       /auth/wechat/status [post]
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

// @Summary      用户
// @Description  查询微信绑定扫码状态
// @ID           wechat-status
// @Tags         user
// @param        code          formData  string  true   "查询code"
// @Success      200
// @Success      302
// @Failure      400  {object}  errs.JsonRespondError
// @Failure      404   {object}  errs.JsonRespondError
// @Router       /auth/bind/wechat/status [post]
func (a *Auth) wechatBindStatus(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		code := ctx.PostForm("code")
		ok, err := a.oauthSvc.WechatScanBindStatus(code)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		return errs2.ErrRecordNotFound
	})
}

func (a *Auth) bbsOAuthCallback(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, ok := token.UserId(ctx)
		code := ctx.Query("code")
		if code == "" {
			return errs.NewBadRequestError(1001, "code不能为空")
		}
		if ok {
			err := a.oauthSvc.BindBbs(uid, code)
			if err != nil {
				return err
			}
			if uri := ctx.Query("redirect"); uri != "" {
				ctx.Redirect(http.StatusFound, uri)
				return nil
			}
			return "绑定成功"
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
			switch code[0] {
			case "login":
				if err := a.oauthSvc.WechatScanLogin(string(msg.FromUserName), code[1]); err != nil {
					logrus.Errorf("wx login message handler: %v", err)
					return nil
				}
			case "bind":
				if err := a.oauthSvc.WechatScanBind(string(msg.FromUserName), code[1]); err != nil {
					logrus.Errorf("wx bind message handler: %v", err)
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
		if err != nil {
			logrus.Errorf("wx proxy: %v", err)
			return
		}
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

func (a *Auth) oauthHandle(ctx *gin.Context, resp *vo.OAuthRespond) interface{} {
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
	if uri := ctx.Query("redirect"); uri != "" {
		ctx.Redirect(http.StatusFound, uri)
		return nil
	}
	return gin.H{
		"uid":   resp.UserInfo.ID,
		"token": tokenString,
	}
}

// @Summary      忘记密码
// @Description  往邮箱里发送一个找回密码的链接
// @ID           forget-password
// @Tags         user
// @Param        email  formData  string  true  "邮箱地址"
// @Success      200
// @Failure      400  {object}  errs.JsonRespondError
// @Failure      404  {object}  errs.JsonRespondError
// @Router       /account/forgot-password [post]
func (a *Auth) forgetPassword(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		return a.safe.Rate(&dto2.SafeUserinfo{
			IP: c.ClientIP(),
		}, &dto2.SafeRule{
			Name:        "forget-password",
			Description: "忘记密码",
			PeriodCnt:   5,
			Period:      3 * time.Minute,
		}, func() error {
			err := a.RequestForgetPasswordEmail(c.PostForm("email"))
			if err != nil {
				return err
			}
			return nil
		})
	})
}

// @Summary      校验重置密码
// @Description  校验重置密码的code
// @ID           valid-reset-password
// @Tags         user
// @Param        code  query     string  true  "重置code"
// @Success      200   {object}  vo.UserInfo
// @Failure      400  {object}  errs.JsonRespondError
// @Failure      404  {object}  errs.JsonRespondError
// @Router       /account/reset-password [get]
func (a *Auth) validResetPassword(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		user, err := a.User.ValidResetPassword(c.Query("code"))
		if err != nil {
			return err
		}
		return user
	})
}

// @Summary      重置密码
// @Description  通过忘记密码的邮件重置密码
// @ID           reset-password
// @Tags         user
// @Param        code        formData  string  true  "重置code"
// @Param        password    formData  string  true  "输入密码"
// @Param        repassword  formData  string  true  "再输入一次密码"
// @Success      200
// @Failure      400  {object}  errs.JsonRespondError
// @Failure      404  {object}  errs.JsonRespondError
// @Router       /account/reset-password [post]
func (a *Auth) resetPassword(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		req := &request.ResetPasswordRequest{}
		err := c.ShouldBind(req)
		if err != nil {
			return err
		}
		return a.User.ResetPassword(req.Code, req.Password)
	})
}

func (a *Auth) Register(r *gin.RouterGroup) {
	rg := r.Group("/account")
	rg.POST("/login", a.login)
	rg.GET("/logout", token.UserAuth(false), a.logout)
	rg.POST("/register", a.register)
	rg.POST("/register/request-email-code", a.requestEmailCode)
	rg.POST("/forgot-password", a.forgetPassword)
	rg.GET("/reset-password", a.validResetPassword)
	rg.POST("/reset-password", a.resetPassword)

	rg = r.Group("/auth")
	rg.GET("/bbs", a.bbsOAuth)
	rg.GET("/bbs/callback", token.UserAuth(false), a.bbsOAuthCallback)
	rg.POST("/wechat", a.wechatOAuth)
	//rg.GET("/wechat/callback", s.wechatOAuthCallback)
	rg.POST("/wechat/request", a.wechatRequest)
	rg.POST("/wechat/status", a.wechatStatus)
	rg.Any("/wechat/handle", a.wechatHandle)

	rg = r.Group("/auth/bind", token.UserAuth(true))
	rg.POST("/wechat/request", a.wechatBindRequest)
	rg.POST("/wechat/status", a.wechatBindStatus)
}
