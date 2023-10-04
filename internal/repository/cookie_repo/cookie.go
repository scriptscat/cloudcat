package cookie_repo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/scriptscat/cloudcat/internal/model/entity/cookie_entity"
	"github.com/scriptscat/cloudcat/pkg/bbolt"
	bolt "go.etcd.io/bbolt"
)

type CookieRepo interface {
	Find(ctx context.Context, storageName, host string) (*cookie_entity.Cookie, error)
	FindPage(ctx context.Context, storageName string) ([]*cookie_entity.Cookie, int64, error)
	Create(ctx context.Context, cookie *cookie_entity.Cookie) error
	Update(ctx context.Context, cookie *cookie_entity.Cookie) error
	Delete(ctx context.Context, storageName, host string) error

	DeleteByStorage(ctx context.Context, storageName string) error
}

var defaultCookie CookieRepo

func Cookie() CookieRepo {
	return defaultCookie
}

func RegisterCookie(i CookieRepo) {
	defaultCookie = i
}

type cookieRepo struct {
}

func NewCookie() CookieRepo {
	return &cookieRepo{}
}

func (u *cookieRepo) Find(ctx context.Context, storageName, host string) (*cookie_entity.Cookie, error) {
	ret := &cookie_entity.Cookie{}
	if err := bbolt.Default().Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("cookie")).CreateBucketIfNotExists([]byte(storageName))
		if err != nil {
			return err
		}
		data := b.Get([]byte(host))
		if data == nil {
			return bbolt.ErrNil
		}
		return json.Unmarshal(data, ret)
	}); err != nil {
		if bbolt.IsNil(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *cookieRepo) Create(ctx context.Context, cookie *cookie_entity.Cookie) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("cookie")).CreateBucketIfNotExists([]byte(cookie.StorageName))
		if err != nil {
			return err
		}
		data, err := json.Marshal(cookie)
		if err != nil {
			return err
		}
		return b.Put([]byte(cookie.Host), data)
	})
}

func (u *cookieRepo) Update(ctx context.Context, cookie *cookie_entity.Cookie) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("cookie")).CreateBucketIfNotExists([]byte(cookie.StorageName))
		if err != nil {
			return err
		}
		data, err := json.Marshal(cookie)
		if err != nil {
			return err
		}
		return b.Put([]byte(cookie.Host), data)
	})
}

func (u *cookieRepo) Delete(ctx context.Context, storageName, host string) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("cookie")).CreateBucketIfNotExists([]byte(storageName))
		if err != nil {
			return err
		}
		return b.Delete([]byte(host))
	})
}

func (u *cookieRepo) FindPage(ctx context.Context, storageName string) ([]*cookie_entity.Cookie, int64, error) {
	var list []*cookie_entity.Cookie
	if err := bbolt.Default().Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte("cookie")).CreateBucketIfNotExists([]byte(storageName))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var cookie cookie_entity.Cookie
			if err := json.Unmarshal(v, &cookie); err != nil {
				return err
			}
			list = append(list, &cookie)
		}
		return nil
	}); err != nil {
		return nil, 0, err
	}
	return list, int64(len(list)), nil
}

func (u *cookieRepo) DeleteByStorage(ctx context.Context, storageName string) error {
	return bbolt.Default().Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("cookie")).DeleteBucket([]byte(storageName))
		if errors.Is(err, bolt.ErrBucketNotFound) {
			return nil
		}
		return err
	})
}
