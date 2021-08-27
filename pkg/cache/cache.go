package cache

import (
	"time"
)

type Cache interface {
	GetOrSet(key string, get interface{}, set func() (interface{}, error), opts ...Option) error
	Set(key string, val interface{}, opts ...Option) error
	Get(key string, get interface{}, opts ...Option) error
	Has(key string) (bool, error)
}

type Depend interface {
	Val() interface{}
	Ok() error
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
