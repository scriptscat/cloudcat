package errs

import (
	"net/http"

	"github.com/scriptscat/cloudcat/pkg/errs"
)

var (
	ErrOAuthPlatformNotSupport    = errs.NewBadRequestError(1000, "三方登录平台不支持")
	ErrOAuthPlatformNotConfigured = errs.NewBadRequestError(1001, "三方登录平台未配置")

	ErrOAuthMustHaveEmail  = errs.NewBadRequestError(1002, "三方登录必须有配置Email")
	ErrRegisterVerifyEmail = errs.NewBadRequestError(1003, "注册必须验证Email")
	ErrEmailVCodeNotFound  = errs.NewBadRequestError(1004, "邮箱验证码不存在")

	ErrUserNotFound        = errs.NewError(http.StatusNotFound, 1004, "没有这个用户")
	ErrEmailSuffixNotAllow = errs.NewBadRequestError(1005, "邮箱后缀不允许")
	ErrRegisterDisable     = errs.NewBadRequestError(1006, "不允许注册")
	ErrRegisterNeedInvCode = errs.NewBadRequestError(1007, "注册需要邀请码")

	ErrWrongPassword  = errs.NewBadRequestError(1008, "登录密码错误")
	ErrEmailExist     = errs.NewBadRequestError(1009, "邮箱已经注册过了")
	ErrMobileExist    = errs.NewBadRequestError(1010, "手机号码已经注册过了")
	ErrUsernameExist  = errs.NewBadRequestError(1011, "用户名已经注册过了")
	ErrRecordNotFound = errs.NewError(http.StatusNotFound, 1012, "记录未找到")

	ErrAvatarNotImage = errs.NewBadRequestError(1013, "上传的头像不是一个正确的图片")
	ErrAvatarIsNil    = errs.NewBadRequestError(1014, "头像是空的")
	ErrAvatarTooBig   = errs.NewBadRequestError(1015, "头像不能超过1M")
	ErrNotUnbind      = errs.NewBadRequestError(10016, "绑定未超过30天,禁止解绑")
	ErrBindOtherUser  = errs.NewBadRequestError(10017, "绑定过其他账号了")
	ErrBindOtherOAuth = errs.NewBadRequestError(10017, "绑定过其他三方账号了")
)
