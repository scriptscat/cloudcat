package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	"gorm.io/gorm"
)

func T1631263155() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1631263155",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&entity.User{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&entity.User{})
		},
	}
}
