package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	entity2 "github.com/scriptscat/cloudcat/internal/service/sync/domain/entity"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

func T1631861289() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1631861289",
		Migrate: func(db *gorm.DB) error {
			return utils.Errs(
				db.AutoMigrate(&entity2.SyncScript{}),
				db.AutoMigrate(&entity2.SyncSubscribe{}),
				db.AutoMigrate(&entity2.SyncDevice{}),
			)
		},
		Rollback: func(db *gorm.DB) error {
			return utils.Errs(
				db.Migrator().DropTable(&entity2.SyncSubscribe{}),
				db.Migrator().DropTable(&entity2.SyncScript{}),
				db.Migrator().DropTable(&entity2.SyncDevice{}),
			)
		},
	}
}
