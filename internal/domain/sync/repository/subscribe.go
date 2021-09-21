package repository

import (
	"crypto/sha1"
	"fmt"

	"github.com/scriptscat/cloudcat/internal/domain/sync/entity"
	"gorm.io/gorm"
)

type Subscribe interface {
	Save(entity *entity.SyncSubscribe) error
	FindByUrl(user, device int64, url string) (*entity.SyncSubscribe, error)
}

type subscribe struct {
	db *gorm.DB
}

func NewSubscribe(db *gorm.DB) Subscribe {
	return &subscribe{db: db}
}

func (s *subscribe) Save(entity *entity.SyncSubscribe) error {
	return s.db.Save(entity).Error
}

func (s *subscribe) FindByUrl(user, device int64, url string) (*entity.SyncSubscribe, error) {
	ret := &entity.SyncSubscribe{}
	if err := s.db.Where("user_id=? and device_id=? and url_hash=?", user, device, fmt.Sprintf("%x", sha1.Sum([]byte(url)))).First(&ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
