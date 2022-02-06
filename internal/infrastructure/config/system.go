package config

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/cloudcat/internal/pkg/kvdb"
)

//go:generate mockgen -source ./system.go -destination ./mock/system.go
type SystemConfig interface {
	GetConfig(key string) (string, error)
	SetConfig(key, value string) error
	DelConfig(key string) error
}

type systemConfig struct {
	kv kvdb.KvDb
}

func NewSystemConfig(kv kvdb.KvDb) SystemConfig {
	return &systemConfig{kv: kv}
}

func (s *systemConfig) GetConfig(key string) (string, error) {
	key = s.key(key)
	v, err := s.kv.Get(context.Background(), key)
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return v, nil
}

func (s *systemConfig) SetConfig(key, value string) error {
	key = s.key(key)
	if err := s.kv.Set(context.Background(), key, value, 0); err != nil {
		return err
	}
	return nil
}

func (s *systemConfig) DelConfig(key string) error {
	if err := s.kv.Del(context.Background(), key); err != nil {
		return err
	}
	return nil
}

func (s *systemConfig) key(key string) string {
	return "system:config:" + key
}
