package api

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/cloudcat/internal/infrastructure/middleware/token"
	service3 "github.com/scriptscat/cloudcat/internal/infrastructure/sender"
	service2 "github.com/scriptscat/cloudcat/internal/service/safe/application"
	dto2 "github.com/scriptscat/cloudcat/internal/service/safe/domain/dto"
	application2 "github.com/scriptscat/cloudcat/internal/service/user/application"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/vo"
	"github.com/scriptscat/cloudcat/pkg/errs"
	"github.com/scriptscat/cloudcat/pkg/httputils"
)

const TokenAuthMaxAge = 432000
const TokenAutoRegen = 259200

type User struct {
	application2.User
	oauthSvc application2.OAuth
	sender   service3.Sender
	safe     service2.Safe
}

func NewUser(svc application2.User, oauthSvc application2.OAuth, safe service2.Safe, sender service3.Sender) *User {
	return &User{User: svc, safe: safe, sender: sender, oauthSvc: oauthSvc}
}

// @Summary      用户
// @Description  用户信息
// @ID           user-info
// @Tags         user
// @Security     BearerAuth
// @Success      200  {object}  vo.UserInfo
// @Failure      403
// @Router       /user [get]
func (u *User) get(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, _ := token.UserId(ctx)
		ret, err := u.UserInfo(uid)
		if err != nil {
			return err
		}
		return ret
	})
}

// @Summary      用户
// @Description  用户设置信息
// @ID           user-setting-info
// @Tags         user
// @Security     BearerAuth
// @Success      200
// @Failure      403
// @Router       /user/setting [get]
func (u *User) setting(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, _ := token.UserId(ctx)
		info, err := u.UserInfo(uid)
		if err != nil {
			return err
		}
		open, err := u.oauthSvc.OAuthPlatform(uid)
		return gin.H{
			"info": info,
			"open": open,
		}
	})
}

// @Summary      用户
// @Description  请求邮箱修改验证码
// @ID           change-email-code
// @Tags         user
// @Security     BearerAuth
// @Param        email  formData  string  true  "邮箱"
// @Success      200
// @Failure      403
// @Router       /user/request-change-email-code [put]
func (u *User) requestChangeEmailCode(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		email := ctx.PostForm("email")
		code, err := u.RequestEmailCode(email, "change-user-info")
		if err != nil {
			return err
		}
		return u.sender.SendEmail(email, "注册验证码", "您的验证码为:"+code.Code+" 请于5分钟内输入", "text/html")
	})
}

// @Summary      用户
// @Description  修改用户信息
// @ID           user-update-info
// @Tags         user
// @Security     BearerAuth
// @Param        username  formData  string  true  "用户名"
// @Param        email     formData  string  true  "邮箱"
// @Param        code      formData  string  true  "邮箱验证码"
// @Success      200
// @Failure      403
// @Router       /user [put]
func (u *User) update(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, _ := token.UserId(ctx)
		req := &vo.UpdateUserInfo{}
		if err := ctx.ShouldBind(req); err != nil {
			return err
		}
		return u.UpdateUserInfo(uid, req)
	})
}

// @Summary      用户
// @Description  修改用户密码
// @ID           user-update-password
// @Security     BearerAuth
// @Tags         user
// @Success      200
// @Failure      403
// @Router       /user/password [put]
func (u *User) password(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, _ := token.UserId(ctx)
		req := &vo.UpdatePassword{}
		if err := ctx.ShouldBind(req); err != nil {
			return err
		}
		return u.UpdatePassword(uid, req)
	})
}

// @Summary      用户
// @Description  当前用户头像
// @ID           user-avatar
// @Tags         user
// @Security     BearerAuth
// @Success      200
// @Failure      403
// @Router       /user/avatar [get]
func (u *User) avatar(ctx *gin.Context) {
	uid, _ := token.UserId(ctx)
	b, err := u.Avatar(uid)
	if err != nil {
		httputils.HandleError(ctx, err)
		return
	}
	ctx.Header("content-type", http.DetectContentType(b))
	ctx.Writer.Write(b)
}

// @Summary      用户
// @Description  更新用户头像
// @ID           user-update-avatar
// @Tags         user
// @Accept       mpfd
// @Security     BearerAuth
// @Param        avatar  formData  file  true  "头像"
// @Success      200
// @Failure      403
// @Router       /user/avatar [put]
func (u *User) updateAvatar(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, _ := token.UserId(ctx)
		file, err := ctx.FormFile("avatar")
		if err != nil {
			return err
		}
		if file.Size > 1024*1024 {
			return errs.NewBadRequestError(1000, "上传的头像过大")
		}
		f, err := file.Open()
		if err != nil {
			return err
		}
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		return u.safe.Rate(&dto2.SafeUserinfo{
			Identifier: strconv.FormatInt(uid, 10),
		}, &dto2.SafeRule{
			Name:        "user-avatar",
			Description: "头像修改过于频繁",
			PeriodCnt:   5,
			Period:      300 * time.Second,
		}, func() error {
			return u.UploadAvatar(uid, b)
		})
	})
}

// @Summary      用户
// @Description  解绑三方登录
// @ID           user-delete-oauth
// @Tags         user
// @Security     BearerAuth
// @Param        platform  formData  string  true  "普通:bbs|wechat"
// @Success      200
// @Failure      403
// @Router       /user/oauth [delete]
func (u *User) deleteOAuth(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, _ := token.UserId(ctx)
		platform := ctx.PostForm("platform")
		return u.oauthSvc.Unbind(uid, platform)
	})
}

func (u *User) Register(r *gin.RouterGroup) {
	rg := r.Group("/user", token.UserAuth(true))
	rg.GET("", u.get)
	rg.PUT("", u.update)
	rg.GET("/setting", u.setting)
	rg.POST("/request-change-email-code", u.requestChangeEmailCode)
	rg.PUT("/password", u.password)
	rg.GET("/avatar", u.avatar)
	rg.PUT("/avatar", u.updateAvatar)
	rg.DELETE("/oauth", u.deleteOAuth)
}
