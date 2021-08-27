package cache

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	goRedis "github.com/go-redis/redis/v8"
)

type redisCache struct {
	redis *goRedis.Client
}

func NewRedisCache(redis *goRedis.Client) *redisCache {
	return &redisCache{
		redis: redis,
	}
}

func (r *redisCache) GetOrSet(key string, get interface{}, set func() (interface{}, error), opts ...Option) error {
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

func (r *redisCache) Get(key string, get interface{}, opts ...Option) error {
	val, err := r.redis.Get(context.Background(), key).Result()
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

func (r *redisCache) Set(key string, val interface{}, opts ...Option) error {
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
	if err := r.redis.Set(context.Background(), key, b, ttl).Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Has(key string) (bool, error) {
	ok, err := r.redis.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	if ok == 1 {
		return true, nil
	}
	return false, nil
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
