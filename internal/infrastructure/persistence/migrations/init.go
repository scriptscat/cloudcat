package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/cloudcat/internal/pkg/database"
)

func RunMigrations(db *database.Database) error {
	return run(db,
		T1631263156,
		T1631861289,
	)
}

func run(db *database.Database, fs ...func() *gormigrate.Migration) error {
	var ms []*gormigrate.Migration
	for _, f := range fs {
		ms = append(ms, f())
	}
	m := gormigrate.New(db.DB, &gormigrate.Options{
		TableName:                 "migrations",
		IDColumnName:              "id",
		IDColumnSize:              200,
		UseTransaction:            true,
		ValidateUnknownMigrations: true,
	}, ms)
	if err := m.Migrate(); err != nil {
		return err
	}
	return nil
}
