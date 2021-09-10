package errs

import (
	"github.com/scriptscat/cloudcat/internal/pkg/errs"
)

var (
	ErrOAuthPlatformNotSupport    = errs.NewBadRequestError(1000, "三方登录平台不支持")
	ErrOAuthPlatformNotConfigured = errs.NewBadRequestError(1001, "三方登录平台未配置")
)
