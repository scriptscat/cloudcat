package errs

import (
	"net/http"

	"github.com/scriptscat/cloudcat/pkg/errs"
)

var (
	ErrOAuthPlatformNotSupport    = errs.NewBadRequestError(10000, "三方登录平台不支持")
	ErrOAuthPlatformNotConfigured = errs.NewBadRequestError(10001, "三方登录平台未配置")

	ErrOAuthMustHaveEmail  = errs.NewBadRequestError(10002, "三方登录必须有配置Email")
	ErrRegisterVerifyEmail = errs.NewBadRequestError(10003, "注册必须验证Email")
	ErrEmailVCodeNotFound  = errs.NewBadRequestError(10004, "邮箱验证码不存在")

	ErrUserNotFound        = errs.NewError(http.StatusNotFound, 10004, "没有这个用户")
	ErrEmailSuffixNotAllow = errs.NewBadRequestError(10005, "邮箱后缀不允许")
	ErrRegisterDisable     = errs.NewBadRequestError(10006, "不允许注册")
	ErrRegisterNeedInvCode = errs.NewBadRequestError(10007, "注册需要邀请码")

	ErrWrongPassword  = errs.NewBadRequestError(10008, "登录密码错误")
	ErrEmailExist     = errs.NewBadRequestError(10009, "邮箱已经注册过了")
	ErrMobileExist    = errs.NewBadRequestError(10010, "手机号码已经注册过了")
	ErrUsernameExist  = errs.NewBadRequestError(10011, "用户名已经注册过了")
	ErrRecordNotFound = errs.NewError(http.StatusNotFound, 10012, "记录未找到")

	ErrAvatarNotImage = errs.NewBadRequestError(10013, "上传的头像不是一个正确的图片")
	ErrAvatarIsNil    = errs.NewBadRequestError(10014, "头像是空的")
	ErrAvatarTooBig   = errs.NewBadRequestError(10015, "头像不能超过1M")
	ErrNotUnbind      = errs.NewBadRequestError(10016, "绑定未超过30天,禁止解绑")
	ErrBindOtherUser  = errs.NewBadRequestError(10017, "绑定过其他账号了")
	ErrBindOtherOAuth = errs.NewBadRequestError(10017, "绑定过其他三方账号了")

	ErrOpenidNotFound = errs.NewBadRequestError(10018, "未找到用户")
)
