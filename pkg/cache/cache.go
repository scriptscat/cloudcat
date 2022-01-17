package cache

import (
	"time"

	kvdb2 "github.com/scriptscat/cloudcat/internal/pkg/kvdb"
)

type Cache interface {
	GetOrSet(key string, get interface{}, set func() (interface{}, error), opts ...Option) error
	Set(key string, val interface{}, opts ...Option) error
	Get(key string, get interface{}, opts ...Option) error
	Has(key string) (bool, error)
	Del(key string) error
}

type Depend interface {
	Val() interface{}
	Ok() error
}

func NewCache(cfg *Config) (Cache, error) {
	var ret Cache
	switch cfg.Type {
	case "redis", "sqlite":
		kvdb, err := kvdb2.NewKvDb(&kvdb2.Config{
			Type: cfg.Type,
			Redis: struct {
				Addr   string
				Passwd string
				DB     int
			}{
				Addr:   cfg.Redis.Addr,
				Passwd: cfg.Redis.Passwd,
				DB:     cfg.Redis.DB,
			},
			Sqlite: struct {
				File string
			}{
				File: cfg.Sqlite.File,
			},
		})
		if err != nil {
			return nil, err
		}
		ret = NewKvdb(kvdb)
	}
	return ret, nil
}

type Option func(*Options)

type Options struct {
	TTL    time.Duration
	Depend Depend
}

func NewOptions(opts ...Option) *Options {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	return options
}

func WithTTL(t time.Duration) Option {
	return func(options *Options) {
		options.TTL = t
	}
}

func WithDepend(depend Depend) Option {
	return func(options *Options) {
		options.Depend = depend
	}
}

type data struct {
	Depend interface{} `json:"depend"`
	Value  interface{} `json:"value"`
}

type StringCache struct {
	String string
}

type IntCache struct {
	Int int
}

type Int64Cache struct {
	Int64 int64
}
