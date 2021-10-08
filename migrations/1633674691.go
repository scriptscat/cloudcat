package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/cloudcat/internal/domain/sync/entity"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

func T1633674691() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1633674691",
		Migrate: func(db *gorm.DB) error {
			return utils.Errs(
				db.Migrator().DropIndex(&entity.SyncSubscribe{}, "device_url"),
				db.Migrator().AddColumn(&entity.SyncSubscribe{}, "url_hash"),
				db.Migrator().CreateIndex(&entity.SyncSubscribe{}, "device_url"),
			)
		},
		Rollback: func(db *gorm.DB) error {
			return utils.Errs(
				db.Migrator().DropIndex(&entity.SyncSubscribe{}, "device_url"),
				db.Migrator().DropColumn(&entity.SyncSubscribe{}, "url_hash"),
				db.Migrator().CreateIndex(&entity.SyncSubscribe{}, "device_url"),
			)
		},
	}
}
