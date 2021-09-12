package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

func T1631263155() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1631263155",
		Migrate: func(db *gorm.DB) error {
			return utils.Errs(
				db.AutoMigrate(&entity.User{}),
				db.AutoMigrate(&entity.WechatOauthUser{}),
				db.AutoMigrate(&entity.BbsOauthUser{}),
				func() error {
					user := &entity.User{
						Username:   "admin",
						Email:      "admin@admin.com",
						Role:       "admin",
						Createtime: time.Now().Unix(),
					}
					if err := user.SetPassword("admin"); err != nil {
						return err
					}
					return db.Save(user).Error
				}(),
			)
		},
		Rollback: func(db *gorm.DB) error {
			return utils.Errs(
				db.Migrator().DropTable(&entity.User{}),
				db.Migrator().DropTable(&entity.WechatOauthUser{}),
				db.Migrator().DropTable(&entity.BbsOauthUser{}),
			)
		},
	}
}
