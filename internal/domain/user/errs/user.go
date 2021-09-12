package errs

import (
	"net/http"

	"github.com/scriptscat/cloudcat/internal/pkg/errs"
)

var (
	ErrOAuthPlatformNotSupport    = errs.NewBadRequestError(1000, "三方登录平台不支持")
	ErrOAuthPlatformNotConfigured = errs.NewBadRequestError(1001, "三方登录平台未配置")

	ErrOAuthMustHaveEmail  = errs.NewBadRequestError(1002, "三方登录必须有配置Email")
	ErrRegisterVerifyEmail = errs.NewBadRequestError(1003, "注册必须验证Email")

	ErrUserNotFound        = errs.NewError(http.StatusNotFound, 1004, "没有这个用户")
	ErrEmailSuffixNotAllow = errs.NewBadRequestError(1005, "邮箱后缀不允许")
	ErrRegisterDisable     = errs.NewBadRequestError(1006, "不允许注册")
	ErrRegisterNeedInvCode = errs.NewBadRequestError(1007, "注册需要邀请码")

	ErrWrongPassword = errs.NewBadRequestError(1008, "登录密码错误")
	ErrEmailExist    = errs.NewBadRequestError(1009, "邮箱已经注册过了")
)
