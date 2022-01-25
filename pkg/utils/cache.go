package utils

import (
	"context"
	"time"

	"github.com/scriptscat/cloudcat/internal/pkg/kvdb"
	"github.com/silenceper/wechat/v2/cache"
)

type WxCache struct {
	kv kvdb.KvDb
}

func NewWxCache(kv kvdb.KvDb) cache.Cache {
	return &WxCache{kv: kv}
}

func (w *WxCache) Get(key string) interface{} {
	ret, err := w.kv.Get(context.Background(), key)
	if err != nil {
		return nil
	}
	return ret
}

func (w *WxCache) Set(key string, val interface{}, timeout time.Duration) error {
	if timeout == 0 {
		timeout = time.Hour * 24 * 30
	}
	return w.kv.Set(context.Background(), key, val.(string), timeout)
}

func (w *WxCache) IsExist(key string) bool {
	if ret, err := w.kv.Has(context.Background(), key); err != nil {
		return false
	} else {
		return ret
	}
}

func (w *WxCache) Delete(key string) error {
	return w.kv.Del(context.Background(), key)
}
