package kvdb

import (
	"context"
	"time"

	goRedis "github.com/go-redis/redis/v8"
	sqliteOrm "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type KvDb interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Has(ctx context.Context, key string) (bool, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Client() interface{}
	DbType() string
}

func NewKvDb(cfg *Config) (KvDb, error) {
	var ret KvDb
	switch cfg.Type {
	case "redis":
		ret = newRedis(goRedis.NewClient(&goRedis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Passwd,
			DB:       cfg.Redis.DB,
		}))
	case "sqlite":
		db, err := gorm.Open(sqliteOrm.Open(cfg.Sqlite.File))
		if err != nil {
			return nil, err
		}
		ret, err = newSqlite(db)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}
