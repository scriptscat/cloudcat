package config

import (
	"context"

	"github.com/scriptscat/cloudcat/pkg/kvdb"
)

type SystemConfig interface {
	GetConfig(key string) (string, error)
	SetConfig(key, value string) error
}

type systemConfig struct {
	kv kvdb.KvDb
}

func NewSystemConfig(kv kvdb.KvDb) SystemConfig {
	return &systemConfig{kv: kv}
}

func (s *systemConfig) GetConfig(key string) (string, error) {
	return s.kv.Get(context.Background(), s.key(key))
}

func (s *systemConfig) SetConfig(key, value string) error {
	return s.kv.Set(context.Background(), key, value, 0)
}

func (s *systemConfig) key(key string) string {
	return "system:config:" + key
}
