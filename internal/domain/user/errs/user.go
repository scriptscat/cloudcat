package errs

import (
	"github.com/scriptscat/cloudcat/internal/pkg/errs"
)

var (
	ErrOAuthPlatformNotSupport    = errs.NewBadRequestError(1000, "三方登录平台不支持")
	ErrOAuthPlatformNotConfigured = errs.NewBadRequestError(1001, "三方登录平台未配置")

	ErrOAuthMustHaveEmail = errs.NewBadRequestError(1002, "三方登录必须有配置Email")
	ErrMustVerifyEmail    = errs.NewBadRequestError(1003, "必须验证Email")
)
