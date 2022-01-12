package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	entity2 "github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

func T1631263155() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1631263155",
		Migrate: func(db *gorm.DB) error {
			return utils.Errs(
				db.AutoMigrate(&entity2.User{}),
				db.AutoMigrate(&entity2.WechatOauthUser{}),
				db.AutoMigrate(&entity2.BbsOauthUser{}),
				func() error {
					user := &entity2.User{
						Username:   "admin",
						Email:      "admin@admin.com",
						Role:       "admin",
						Createtime: time.Now().Unix(),
						Updatetime: 0,
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
				db.Migrator().DropTable(&entity2.User{}),
				db.Migrator().DropTable(&entity2.WechatOauthUser{}),
				db.Migrator().DropTable(&entity2.BbsOauthUser{}),
			)
		},
	}
}
