package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

func T1631263156() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1631263156",
		Migrate: func(db *gorm.DB) error {
			return utils.Errs(
				db.AutoMigrate(&entity.User{}),
				db.AutoMigrate(&entity.VerifyCode{}),
				db.AutoMigrate(&entity.BbsOauthUser{}),
				db.AutoMigrate(&entity.WechatOauthUser{}),
			)
		},
		Rollback: func(db *gorm.DB) error {
			return utils.Errs(
				db.Migrator().DropTable(&entity.User{}),
				db.Migrator().DropTable(&entity.VerifyCode{}),
				db.Migrator().DropTable(&entity.BbsOauthUser{}),
				db.Migrator().DropTable(&entity.WechatOauthUser{}),
			)
		},
	}
}
