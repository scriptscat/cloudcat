package cache

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/scriptscat/cloudcat/internal/infrastructure/kvdb"
)

type kvdbCache struct {
	kv kvdb.KvDb
}

func NewKvdb(kv kvdb.KvDb) *kvdbCache {
	return &kvdbCache{
		kv: kv,
	}
}

func (r *kvdbCache) GetOrSet(key string, get interface{}, set func() (interface{}, error), opts ...Option) error {
	err := r.Get(key, get, opts...)
	if err != nil {
		val, err := set()
		if err != nil {
			return err
		}
		if err := r.Set(key, val, opts...); err != nil {
			return err
		}
	}
	return nil
}

func (r *kvdbCache) Get(key string, get interface{}, opts ...Option) error {
	val, err := r.kv.Get(context.Background(), key)
	if err != nil {
		return err
	}
	options := NewOptions(opts...)
	ret := &data{Value: get, Depend: options.Depend}
	if err := json.Unmarshal([]byte(val), ret); err != nil {
		return err
	}
	if options.Depend != nil {
		if err := options.Depend.Ok(); err != nil {
			return err
		}
	}
	return nil
}

func (r *kvdbCache) Set(key string, val interface{}, opts ...Option) error {
	options := NewOptions(opts...)
	ttl := time.Duration(0)
	if options.TTL > 0 {
		ttl = options.TTL
	}
	data := &data{Value: val}
	if options.Depend != nil {
		data.Depend = options.Depend.Val()
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := r.kv.Set(context.Background(), key, string(b), ttl); err != nil {
		return err
	}
	return nil
}

func (r *kvdbCache) Has(key string) (bool, error) {
	ok, err := r.kv.Has(context.Background(), key)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *kvdbCache) Del(key string) error {
	return r.kv.Del(context.Background(), key)
}

func copyInterface(dst interface{}, src interface{}) {
	dstof := reflect.ValueOf(dst)
	if dstof.Kind() == reflect.Ptr {
		el := dstof.Elem()
		srcof := reflect.ValueOf(src)
		if srcof.Kind() == reflect.Ptr {
			el.Set(srcof.Elem())
		} else {
			el.Set(srcof)
		}
	}
}
