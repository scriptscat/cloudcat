package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	"github.com/scriptscat/cloudcat/internal/pkg/cnt"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
	"gorm.io/gorm"
)

type BBSOAuth interface {
	FindByOpenid(openid string) (*entity.BbsOauthUser, error)
	Save(bbs *entity.BbsOauthUser) error
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

func (b *bbsOAuth) Save(bbs *entity.BbsOauthUser) error {
	return b.db.Save(bbs).Error
}

type WechatOAuth interface {
	Save(u *entity.WechatOauthUser) error
	FindByOpenid(openid string) (*entity.WechatOauthUser, error)
	BindCodeUid(code string, uid int64) error
	FindCodeUid(code string) (int64, error)
}

type wechatOAuth struct {
	db *gorm.DB
	kv kvdb.KvDb
}

func NewWechatOAuth(db *gorm.DB, kv kvdb.KvDb) WechatOAuth {
	return &wechatOAuth{
		db: db,
		kv: kv,
	}
}

func (b *wechatOAuth) Save(u *entity.WechatOauthUser) error {
	return b.db.Save(u).Error
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

func (b *wechatOAuth) BindCodeUid(code string, uid int64) error {
	return b.kv.Set(context.Background(), b.key(code), strconv.FormatInt(uid, 10), time.Minute*20)
}

func (b *wechatOAuth) FindCodeUid(code string) (int64, error) {
	result, err := b.kv.Get(context.Background(), b.key(code))
	if err != nil {
		return 0, err
	}
	if result == "" {
		return 0, nil
	}
	if err := b.kv.Del(context.Background(), b.key(code)); err != nil {
		return 0, err
	}
	return strconv.ParseInt(result, 10, 64)
}

func (b *wechatOAuth) key(code string) string {
	return "user:oauth:code:" + code
}
