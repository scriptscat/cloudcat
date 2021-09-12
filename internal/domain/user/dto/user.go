package dto

import "github.com/scriptscat/cloudcat/internal/domain/user/entity"

type UserInfo struct {
	ID         int64  `json:"id"`       // 用户id
	Username   string `json:"username"` // 用户名
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Createtime int64  `json:"createtime"`
	Updatetime int64  `json:"updatetime"`
}

func ToUserInfo(user *entity.User) *UserInfo {
	return &UserInfo{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Mobile:     user.Mobile,
		Createtime: user.Createtime,
		Updatetime: user.Updatetime,
	}
}

type OAuthRespond struct {
	UserInfo *UserInfo
	IsBind   bool
}
