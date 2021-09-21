package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/cloudcat/internal/domain/sync/entity"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
)

// TODO: 考虑value是否需要云同步,或者只同步userconfig
// Setting 同步Value和系统设置
type Setting interface {
	SaveValue(value *entity.SyncValue) error
}

type setting struct {
	sync.RWMutex
	kv    kvdb.KvDb
	redis *redis.Client
}

// NewSetting kvdb必须是redis才能储存value数据,否则是直接返回成功和空
func NewSetting(kv kvdb.KvDb) Setting {
	var rds *redis.Client
	if kv.DbType() == "redis" {
		rds = kv.Client().(*redis.Client)
	}
	return &setting{kv: kv, redis: rds}
}

func (s *setting) SaveValue(value *entity.SyncValue) error {
	rds, err := s.rds()
	if err != nil {
		return err
	}
	hkey := ""
	if value.StorageName == "" {
		hkey = "script:" + value.ScriptUUID + ":" + value.Key
	} else {
		hkey = "storage:" + value.StorageName + ":" + value.Key
	}
	return rds.HSet(context.Background(), "sync:value:"+fmt.Sprintf("%d:%d", value.UserID, value.DeviceID),
		hkey, value.Value).Err()
}

func (s *setting) rds() (*redis.Client, error) {
	if s.redis == nil {
		return nil, errors.New("请使用redis作为kvdb,否则无法存储value数据")
	}
	return s.redis, nil
}
