package kvdb

import (
	"context"
	"strconv"
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
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
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

func (s *sqlite) Del(ctx context.Context, key string) error {
	m := &kvTable{Key: key}
	return s.db.Delete(m).Error
}

func (s *sqlite) IncrBy(ctx context.Context, key string, value int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Exec("BEGIN EXCLUSIVE TRANSACTION;")
		m := &kvTable{Key: key}
		if err := s.db.First(m).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		} else if m.Expired != 0 && time.Now().Unix() > m.Expired {
			m = &kvTable{
				Key:     key,
				Value:   "",
				Expired: 0,
			}
		}
		t, _ := strconv.ParseInt(m.Value, 10, 64)
		m.Value = strconv.FormatInt(t+value, 10)
		return tx.Save(m).Error
	})
}

func (s *sqlite) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return s.db.Model(&kvTable{}).Where("key=?", key).Update("expired", time.Now().Add(expiration).Unix()).Error
}

func (s *sqlite) Client() interface{} {
	panic("implement me")
}

func (s *sqlite) DbType() string {
	return "sqlite"
}
