package dto

import "github.com/scriptscat/cloudcat/internal/domain/user/entity"

type UserInfo struct {
	ID         int64  `json:"id"` // 用户id
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Role       string `json:"role"`
	Createtime int64  `json:"createtime"`
	Updatetime int64  `json:"updatetime"`
}

func ToUserInfo(user *entity.User) *UserInfo {
	return &UserInfo{
		ID:         user.ID,
		Username:   user.Username,
		Avatar:     user.Avatar,
		Role:       user.Role,
		Createtime: user.Createtime,
		Updatetime: user.Updatetime,
	}
}

type OAuthRespond struct {
	UserInfo *UserInfo
	IsBind   bool
}

type Login struct {
	Username  string `form:"username" binding:"required" label:"用户名或邮箱"`
	Password  string `form:"password" binding:"required" label:"密码"`
	AutoLogin bool   `form:"auto_login" label:"自动登录"`
}

type Register struct {
	Username   string `form:"username" binding:"required,min=3,max=16" label:"用户名"`
	Email      string `form:"email" binding:"required,min=3,max=32,email" label:"邮箱"`
	Password   string `form:"password" binding:"required,min=6,max=18" label:"密码"`
	Repassword string `form:"password" binding:"required,min=6,max=18,eqfield=Password" label:"再输入一次密码"`
	// 开启邮箱验证
	EmailVerifyCode string `form:"email_verify_code" binding:"omitempty,len=6,alphanum" label:"邮箱验证码"`
	// 开启邀请码注册
	InvCode string `form:"inv_code" binding:"omitempty,len=6,alphanum" label:"邀请码"`
}

type VerifyEmail struct {
	Code string `json:"code"`
}

type WechatScanLogin struct {
	URL  string `json:"url"`
	Code string `json:"code"`
}
