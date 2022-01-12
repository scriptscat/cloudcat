package repository

import (
	"github.com/scriptscat/cloudcat/internal/service/sync/domain/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -source ./device.go -destination ./mock/device.go

type Device interface {
	FindById(id int64) (*entity.SyncDevice, error)
	ListDevice(user int64) ([]*entity.SyncDevice, error)
	Save(device *entity.SyncDevice) error
	UpdateSetting(id int64, setting string, settingtime int64) error
}

type device struct {
	db *gorm.DB
}

func NewDevice(db *gorm.DB) Device {
	return &device{db: db}
}

func (d *device) FindById(id int64) (*entity.SyncDevice, error) {
	ret := &entity.SyncDevice{ID: id}
	if err := d.db.First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
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

func (d *device) UpdateSetting(id int64, setting string, settingtime int64) error {
	return d.db.Model(&entity.SyncDevice{}).Where("id=?", id).Updates(map[string]interface{}{
		"setting":     setting,
		"settingtime": settingtime,
	}).Error
}
