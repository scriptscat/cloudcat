package repository

import (
	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	"github.com/scriptscat/cloudcat/internal/pkg/cnt"
	"gorm.io/gorm"
)

type BBSOAuth interface {
	FindByOpenid(openid string) (*entity.BbsOauthUser, error)
}

type bbsOAuth struct {
	db *gorm.DB
}

func NewBbsOAuth(db *gorm.DB) BBSOAuth {
	return &bbsOAuth{
		db: db,
	}
}

func (b *bbsOAuth) FindByOpenid(openid string) (*entity.BbsOauthUser, error) {
	ret := &entity.BbsOauthUser{}
	if err := b.db.First(ret, "openid=? and status=?", openid, cnt.ACTIVE).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

type WechatOAuth interface {
	FindByOpenid(openid string) (*entity.WechatOauthUser, error)
}

type wechatOAuth struct {
	db *gorm.DB
}

func NewWechatOAuth(db *gorm.DB) WechatOAuth {
	return &wechatOAuth{
		db: db,
	}
}

func (b *wechatOAuth) FindByOpenid(openid string) (*entity.WechatOauthUser, error) {
	ret := &entity.WechatOauthUser{}
	if err := b.db.First(ret, "openid=? and status=?", openid, cnt.ACTIVE).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
