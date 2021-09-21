package repository

import (
	"github.com/scriptscat/cloudcat/internal/domain/sync/entity"
	"github.com/scriptscat/cloudcat/internal/pkg/cnt"
	"gorm.io/gorm"
)

type Script interface {
	ListScript(user, device int64) ([]*entity.SyncScript, error)
	FindByUUID(user, device int64, uuid string) (*entity.SyncScript, error)
	Save(entity *entity.SyncScript) error
}

type script struct {
	db *gorm.DB
}

func NewScript(db *gorm.DB) Script {
	return &script{
		db: db,
	}
}

func (s *script) ListScript(user, device int64) ([]*entity.SyncScript, error) {
	list := make([]*entity.SyncScript, 0)
	if err := s.db.Model(&entity.SyncScript{}).Where("user_id=? and device_id=? and state=?", user, device, cnt.ACTIVE).Scan(&list).Error; err != nil {
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
