package persistence

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/cloudcat/internal/pkg/kvdb"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/errs"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/repository"
	"github.com/scriptscat/cloudcat/pkg/cnt"
	"gorm.io/gorm"
)

type bbsOAuth struct {
	db *gorm.DB
}

func NewBbsOAuth(db *gorm.DB) repository.BBSOAuth {
	return &bbsOAuth{
		db: db,
	}
}

func (b *bbsOAuth) FindByOpenid(openid string) (*entity.BbsOauthUser, error) {
	ret := &entity.BbsOauthUser{}
	if err := b.db.First(ret, "openid=? and status=?", openid, cnt.ACTIVE).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrOpenidNotFound
		}
		return nil, err
	}
	return ret, nil
}

func (b *bbsOAuth) FindByUid(uid int64) (*entity.BbsOauthUser, error) {
	ret := &entity.BbsOauthUser{}
	if err := b.db.First(ret, "user_id=? and status=?", uid, cnt.ACTIVE).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return ret, nil
}

func (b *bbsOAuth) Save(bbs *entity.BbsOauthUser) error {
	return b.db.Save(bbs).Error
}

func (b *bbsOAuth) Delete(id int64) error {
	return b.db.Delete(&entity.BbsOauthUser{ID: id}).Error
}

type wechatOAuth struct {
	db *gorm.DB
	kv kvdb.KvDb
}

func NewWechatOAuth(db *gorm.DB, kv kvdb.KvDb) repository.WechatOAuth {
	return &wechatOAuth{
		db: db,
		kv: kv,
	}
}

func (w *wechatOAuth) Save(u *entity.WechatOauthUser) error {
	return w.db.Save(u).Error
}

func (w *wechatOAuth) Delete(id int64) error {
	return w.db.Delete(&entity.WechatOauthUser{ID: id}).Error
}

func (w *wechatOAuth) FindByUid(uid int64) (*entity.WechatOauthUser, error) {
	ret := &entity.WechatOauthUser{}
	if err := w.db.First(ret, "user_id=? and status=?", uid, cnt.ACTIVE).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return ret, nil
}

func (w *wechatOAuth) FindByOpenid(openid string) (*entity.WechatOauthUser, error) {
	ret := &entity.WechatOauthUser{}
	if err := w.db.First(ret, "openid=? and status=?", openid, cnt.ACTIVE).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrOpenidNotFound
		}
		return nil, err
	}
	return ret, nil
}

func (w *wechatOAuth) BindCodeUid(code string, uid int64) error {
	return w.kv.Set(context.Background(), w.key(code), strconv.FormatInt(uid, 10), time.Minute*20)
}

func (w *wechatOAuth) FindCodeUid(code string) (int64, error) {
	result, err := w.kv.Get(context.Background(), w.key(code))
	if err != nil {
		if err == redis.Nil {
			return 0, errs.ErrRecordNotFound
		}
		return 0, err
	}
	return strconv.ParseInt(result, 10, 64)
}

func (w *wechatOAuth) DelCode(code string) error {
	return w.kv.Del(context.Background(), w.key(code))
}

func (w *wechatOAuth) key(code string) string {
	return "user:oauth:code:" + code
}
