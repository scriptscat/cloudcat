package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/cloudcat/pkg/database"
)

func RunMigrations(db *database.Database) error {
	return run(db,
		T1631263155,
	)
}

func run(db *database.Database, fs ...func() *gormigrate.Migration) error {
	ms := []*gormigrate.Migration{}
	for _, f := range fs {
		ms = append(ms, f())
	}
	m := gormigrate.New(db.DB, gormigrate.DefaultOptions, ms)
	if err := m.Migrate(); err != nil {
		return err
	}
	return nil
}
