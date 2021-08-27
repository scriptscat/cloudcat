package kvdb

import (
	"context"
	"time"

	goRedis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type kvTable struct {
	Key     string `gorm:"primarykey;column:key;type:varchar(255);index:key,unique"`
	Value   string `gorm:"column:value;type:text"`
	Expired int64  `gorm:"column:expired;type:text"`
}

type sqlite struct {
	db *gorm.DB
}

func newSqlite(db *gorm.DB) (KvDb, error) {
	if err := db.Migrator().AutoMigrate(&kvTable{}); err != nil {
		return nil, err
	}
	return &sqlite{db: db}, nil
}

func (s *sqlite) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m := &kvTable{Key: key, Value: key, Expired: time.Now().Add(expiration).Unix()}
	return s.db.Save(m).Error
}

func (s *sqlite) Get(ctx context.Context, key string) (string, error) {
	m := &kvTable{Key: key}
	if err := s.db.First(m); err != nil {
		return "", goRedis.Nil
	}
	if time.Now().Unix() > m.Expired {
		return "", goRedis.Nil
	}
	return m.Value, nil
}

func (s *sqlite) Has(ctx context.Context, key string) (bool, error) {
	_, err := s.Get(ctx, key)
	if err != nil {
		if err == goRedis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *sqlite) Client() interface{} {
	panic("implement me")
}

func (s *sqlite) DbType() string {
	return "sqlite"
}
