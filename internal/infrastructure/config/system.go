package config

import (
	"context"

	"github.com/scriptscat/cloudcat/internal/infrastructure/kvdb"
)

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
