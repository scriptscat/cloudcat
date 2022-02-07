package config

import (
	"gorm.io/gorm"
)

//go:generate mockgen -source ./system.go -destination ./mock/system.go
type SystemConfig interface {
	GetConfig(key string) (string, error)
	SetConfig(key, value string) error
	DelConfig(key string) error
}

type systemConfig struct {
	db *gorm.DB
}

type SystemConfigTable struct {
	ID    int64  `gorm:"primaryKey" json:"-"`
	Key   string `gorm:"column:key;unique;type:varchar(255)" json:"key" form:"key"`
	Value string `gorm:"column:value;type:text" json:"value" form:"value"`
}

// NewSystemConfig TODO: 增加缓存
func NewSystemConfig(db *gorm.DB) (SystemConfig, error) {
	if err := db.AutoMigrate(&SystemConfigTable{}); err != nil {
		return nil, err
	}
	return &systemConfig{db: db}, nil
}

func (s *systemConfig) GetConfig(key string) (string, error) {
	m := &SystemConfigTable{}
	if err := s.db.First(m, "`key`=?", key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return m.Value, nil
}

func (s *systemConfig) SetConfig(key, value string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		m := &SystemConfigTable{}
		if err := tx.First(m, "`key`=?", key).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		m.Key = key
		m.Value = value
		return tx.Save(m).Error
	})
}

func (s *systemConfig) DelConfig(key string) error {
	return s.db.Delete("`key`=?", key).Error
}
