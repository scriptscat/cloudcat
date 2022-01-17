package persistence

import (
	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/repository"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

type Repositories struct {
	db         *gorm.DB
	User       repository.User
	VerifyCode repository.VerifyCode
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		db:         db,
		User:       NewUser(db),
		VerifyCode: NewVerifyCode(db),
	}
}

func (r *Repositories) AutoMigrate() error {
	return utils.Errs(
		r.db.AutoMigrate(&entity.User{}),
		r.db.AutoMigrate(&entity.VerifyCode{}),
		r.db.AutoMigrate(&entity.BbsOauthUser{}),
		r.db.AutoMigrate(&entity.WechatOauthUser{}),
	)
}
