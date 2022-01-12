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
	"github.com/scriptscat/cloudcat/pkg/cnt"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

//go:generate mockgen -source ./script.go -destination ./mock/script.go

type Script interface {
	LatestVersion(user, device int64) (int64, error)
	PushVersion(user, device int64, data []*dto.SyncScript) (int64, error)
	ActionList(user, device, version int64) ([][]*dto.SyncScript, error)
	ListScript(user, device int64) ([]*entity.SyncScript, error)
	FindByUUID(user, device int64, uuid string) (*entity.SyncScript, error)
	Save(entity *entity.SyncScript) error
	SetStatus(id int64, status int8) error
}

type script struct {
	db    *gorm.DB
	kv    kvdb.KvDb
	redis *redis.Client
}

func NewScript(db *gorm.DB, kv kvdb.KvDb) Script {
	var rds *redis.Client
	if kv.DbType() == "redis" {
		rds = kv.Client().(*redis.Client)
	}
	return &script{
		db:    db,
		kv:    kv,
		redis: rds,
	}
}

func (s *script) LatestVersion(user, device int64) (int64, error) {
	result, err := s.kv.Get(context.Background(), s.key(user, device)+":version")
	if err != nil {
		return 0, err
	}
	return utils.StringToInt64(result), nil
}

func (s *script) PushVersion(user, device int64, data []*dto.SyncScript) (int64, error) {
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

func (s *script) ActionList(user, device, version int64) ([][]*dto.SyncScript, error) {
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
	ret := make([][]*dto.SyncScript, 0)
	for _, v := range list {
		s := make([]*dto.SyncScript, 0)
		if err := json.Unmarshal([]byte(v), &s); err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

func (s *script) key(user, device int64) string {
	return fmt.Sprintf("sync:script:list:%d:%d", user, device)
}

func (s *script) ListScript(user, device int64) ([]*entity.SyncScript, error) {
	list := make([]*entity.SyncScript, 0)
	if err := s.db.Model(&entity.SyncScript{}).Where("user_id=? and device_id=? and status=?", user, device, cnt.ACTIVE).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *script) FindByUUID(user, device int64, uuid string) (*entity.SyncScript, error) {
	ret := &entity.SyncScript{}
	if err := s.db.Where("user_id=? and device_id=? and uuid=?", user, device, uuid).First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *script) Save(entity *entity.SyncScript) error {
	return s.db.Save(entity).Error
}

func (s *script) SetStatus(id int64, status int8) error {
	return s.db.Model(&entity.SyncScript{}).Where("id=?", id).Update("status", status).Error
}

func (s *script) rds() (*redis.Client, error) {
	if s.redis == nil {
		return nil, errors.New("请使用redis作为kvdb,否则无法存储script数据")
	}
	return s.redis, nil
}
