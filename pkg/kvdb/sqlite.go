package kvdb

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type kvTable struct {
	Key     string `gorm:"primarykey;column:key;type:varchar(255);index:key,unique"`
	Value   string `gorm:"column:value;type:text"`
	Expired int64  `gorm:"column:expired;type:bigint;index:expired"`
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

func (s *sqlite) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	m := &kvTable{Key: key, Value: value, Expired: time.Now().Add(expiration).Unix()}
	return s.db.Save(m).Error
}

func (s *sqlite) Get(ctx context.Context, key string) (string, error) {
	m := &kvTable{Key: key}
	if err := s.db.First(m).Error; err != nil {
		return "", err
	}
	if m.Expired != 0 && time.Now().Unix() > m.Expired {
		return "", nil
	}
	return m.Value, nil
}

func (s *sqlite) Has(ctx context.Context, key string) (bool, error) {
	_, err := s.Get(ctx, key)
	if err != nil {
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
