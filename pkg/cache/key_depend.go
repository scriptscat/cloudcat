package cache

import (
	"errors"
	"time"
)

type KeyDepend struct {
	store Cache
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

func NewKeyDepend(store Cache, key string) *KeyDepend {
	return &KeyDepend{
		store: store,
		Key:   key,
	}
}

func WithKeyDepend(store Cache, key string) Option {
	return func(options *Options) {
		options.Depend = NewKeyDepend(store, key)
	}
}

func (v *KeyDepend) InvalidKey() error {
	return v.store.Set(v.Key, &KeyDepend{Key: v.Key, Value: time.Now().Unix()})
}

func (v *KeyDepend) Val() interface{} {
	ret := &KeyDepend{}
	if err := v.store.Get(v.Key, ret); err != nil {
		if err := v.InvalidKey(); err != nil {
			return err
		}
		return &KeyDepend{Key: v.Key, Value: time.Now().Unix()}
	}
	return ret
}

func (v *KeyDepend) Ok() error {
	val := v.Val().(*KeyDepend)
	if v.Value == val.Value {
		return nil
	}
	return errors.New("val not equal")
}
