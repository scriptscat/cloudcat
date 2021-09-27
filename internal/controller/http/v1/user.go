package v1

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	service2 "github.com/scriptscat/cloudcat/internal/domain/safe/service"
	service3 "github.com/scriptscat/cloudcat/internal/domain/system/service"
	"github.com/scriptscat/cloudcat/internal/domain/user/service"
	"github.com/scriptscat/cloudcat/internal/pkg/errs"
	"github.com/scriptscat/cloudcat/internal/pkg/httputils"
)

const TokenAuthMaxAge = 432000
const TokenAutoRegen = 259200

type User struct {
	service.User
	oauthSvc service.OAuth
	sender   service3.Sender
	safe     service2.Safe
}

func NewUser(svc service.User, oauthSvc service.OAuth, safe service2.Safe, sender service3.Sender) *User {
	return &User{User: svc, safe: safe, sender: sender, oauthSvc: oauthSvc}
}

// @Summary     用户
// @Description 用户信息
// @ID          user-info
// @Tags  	    user
// @Success     200 {object} dto.UserInfo
// @Failure     403
// @Router      /user [get]
func (u *User) get(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, _ := userId(ctx)
		ret, err := u.UserInfo(uid)
		if err != nil {
			return err
		}
		return ret
	})
}

// @Summary     用户
// @Description 修改用户信息
// @ID          user-update-info
// @Tags  	    user
// @Success     200
// @Failure     403
// @Router      /user/avatar [put]
func (u *User) update(ctx *gin.Context) {
}

// @Summary     用户
// @Description 当前用户头像
// @ID          user-avatar
// @Tags  	    user
// @Security    BearerAuth
// @Success     200
// @Failure     403
// @Router      /user/avatar [get]
func (u *User) avatar(ctx *gin.Context) {
	uid, _ := userId(ctx)
	b, err := u.Avatar(uid)
	if err != nil {
		httputils.HandleError(ctx, err)
		return
	}
	ctx.Header("content-type", http.DetectContentType(b))
	ctx.Writer.Write(b)
}

// @Summary     用户
// @Description 更新用户头像
// @ID          user-update-avatar
// @Tags  	    user
// @Accept      mpfd
// @Security    BearerAuth
// @Param       avatar formData file true "头像"
// @Success     200
// @Failure     403
// @Router      /user/avatar [put]
func (u *User) updateAvatar(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, _ := userId(ctx)
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
		return u.UploadAvatar(uid, b)
	})
}

func (u *User) Register(r *gin.RouterGroup) {
	rg := r.Group("/user", userAuth())
	rg.GET("", u.get)
	rg.PUT("", u.update)
	rg.GET("/avatar", u.avatar)
	rg.PUT("/avatar", u.updateAvatar)
}
