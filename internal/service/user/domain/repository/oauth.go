package repository

import (
	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
)

//go:generate mockgen -source ./oauth.go -destination ./mock/oauth.go
type BBSOAuth interface {
	FindByOpenid(openid string) (*entity.BbsOauthUser, error)
	FindByUid(uid int64) (*entity.BbsOauthUser, error)
	Save(bbs *entity.BbsOauthUser) error
	Delete(id int64) error
}

type WechatOAuth interface {
	Save(u *entity.WechatOauthUser) error
	FindByOpenid(openid string) (*entity.WechatOauthUser, error)
	FindByUid(uid int64) (*entity.WechatOauthUser, error)
	BindCodeUid(code string, uid int64) error
	FindCodeUid(code string) (int64, error)
	//NOTE:软删除比较好

	Delete(id int64) error
}
