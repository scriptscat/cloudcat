package kvdb

import (
	"context"
	"time"

	goRedis "github.com/go-redis/redis/v8"
)

type redis struct {
	cli *goRedis.Client
}

func newRedis(cli *goRedis.Client) KvDb {
	return &redis{cli: cli}
}

func (r *redis) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.cli.Set(ctx, key, value, expiration).Err()
}

func (r *redis) Get(ctx context.Context, key string) (string, error) {
	ret, err := r.cli.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (r *redis) Del(ctx context.Context, key string) error {
	return r.cli.Del(ctx, key).Err()
}

func (r *redis) Has(ctx context.Context, key string) (bool, error) {
	ok, err := r.cli.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if ok == 1 {
		return true, nil
	}
	return false, nil
}

func (r *redis) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.cli.IncrBy(ctx, key, value).Result()
}

func (r *redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.cli.Expire(ctx, key, expiration).Err()
}

func (r *redis) Client() interface{} {
	return r.cli
}

func (r *redis) DbType() string {
	return "redis"
}
