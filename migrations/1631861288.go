package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/cloudcat/internal/domain/sync/entity"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

func T1631861288() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1631861288",
		Migrate: func(db *gorm.DB) error {
			return utils.Errs(
				db.AutoMigrate(&entity.SyncScript{}),
				db.AutoMigrate(&entity.SyncSubscribe{}),
				db.AutoMigrate(&entity.SyncDevice{}),
			)
		},
		Rollback: func(db *gorm.DB) error {
			return utils.Errs(
				db.Migrator().DropTable(&entity.SyncSubscribe{}),
				db.Migrator().DropTable(&entity.SyncScript{}),
				db.Migrator().DropTable(&entity.SyncDevice{}),
			)
		},
	}
}
