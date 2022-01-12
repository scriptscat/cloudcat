package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/cloudcat/internal/infrastructure/kvdb"
	"github.com/scriptscat/cloudcat/internal/service/sync/domain/dto"
	"github.com/scriptscat/cloudcat/internal/service/sync/domain/entity"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

//go:generate mockgen -source ./subscribe.go -destination ./mock/subscribe.go

type Subscribe interface {
	LatestVersion(user, device int64) (int64, error)
	PushVersion(user, device int64, data []*dto.SyncSubscribe) (int64, error)
	ActionList(user, device, version int64) ([][]*dto.SyncSubscribe, error)
	FindByUrl(user, device int64, url string) (*entity.SyncSubscribe, error)
	Save(entity *entity.SyncSubscribe) error
	SetStatus(id int64, status int8) error
}

type subscribe struct {
	db    *gorm.DB
	kv    kvdb.KvDb
	redis *redis.Client
}

func NewSubscribe(db *gorm.DB, kv kvdb.KvDb) Subscribe {
	var rds *redis.Client
	if kv.DbType() == "redis" {
		rds = kv.Client().(*redis.Client)
	}
	return &subscribe{
		db:    db,
		kv:    kv,
		redis: rds,
	}
}

func (s *subscribe) LatestVersion(user, device int64) (int64, error) {
	result, err := s.kv.Get(context.Background(), s.key(user, device)+":version")
	if err != nil {
		return 0, err
	}
	return utils.StringToInt64(result), nil
}

func (s *subscribe) PushVersion(user, device int64, data []*dto.SyncSubscribe) (int64, error) {
	rds, err := s.rds()
	if err != nil {
		return 0, err
	}
	b, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	version, err := s.kv.IncrBy(context.Background(), s.key(user, device)+":version", 1)
	if err != nil {
		return 0, err
	}
	err = rds.ZAdd(context.Background(), s.key(user, device), &redis.Z{
		Score:  float64(version),
		Member: b,
	}).Err()
	if err != nil {
		return 0, err
	}
	return version, nil
}

func (s *subscribe) ActionList(user, device, version int64) ([][]*dto.SyncSubscribe, error) {
	rds, err := s.rds()
	if err != nil {
		return nil, err
	}
	list, err := rds.ZRangeByScore(context.Background(), s.key(user, device), &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", version),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, err
	}
	ret := make([][]*dto.SyncSubscribe, 0)
	for _, v := range list {
		s := make([]*dto.SyncSubscribe, 0)
		if err := json.Unmarshal([]byte(v), &s); err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

func (s *subscribe) key(user, device int64) string {
	return fmt.Sprintf("sync:subscribe:list:%d:%d", user, device)
}

func (s *subscribe) FindByUrl(user, device int64, url string) (*entity.SyncSubscribe, error) {
	ret := &entity.SyncSubscribe{}
	if err := s.db.Where("user_id=? and device_id=? and url=?", user, device, url).First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *subscribe) Save(entity *entity.SyncSubscribe) error {
	return s.db.Save(entity).Error
}

func (s *subscribe) SetStatus(id int64, status int8) error {
	return s.db.Model(&entity.SyncSubscribe{}).Where("id=?", id).Update("status", status).Error
}

func (s *subscribe) rds() (*redis.Client, error) {
	if s.redis == nil {
		return nil, errors.New("请使用redis作为kvdb,否则无法存储script数据")
	}
	return s.redis, nil
}
