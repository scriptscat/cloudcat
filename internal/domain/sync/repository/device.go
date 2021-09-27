package repository

import (
	"github.com/scriptscat/cloudcat/internal/domain/sync/entity"
	"gorm.io/gorm"
)

type Device interface {
	ListDevice(user int64) ([]*entity.SyncDevice, error)
	Save(device *entity.SyncDevice) error
}

type device struct {
	db *gorm.DB
}

func NewDevice(db *gorm.DB) Device {
	return &device{db: db}
}

func (d *device) ListDevice(user int64) ([]*entity.SyncDevice, error) {
	ret := make([]*entity.SyncDevice, 0)
	if err := d.db.Model(&entity.SyncDevice{}).Where("user_id=?", user).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (d *device) Save(device *entity.SyncDevice) error {
	return d.db.Save(device).Error
}
