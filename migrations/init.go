package migrations

import (
	bbolt2 "github.com/scriptscat/cloudcat/pkg/bbolt"
	"go.etcd.io/bbolt"
)

// RunMigrations 数据库迁移操作
func RunMigrations(db *bbolt.DB) error {
	return run(db,
		T20230210,
	)
}

func run(db *bbolt.DB, fs ...func() *bbolt2.Migration) error {
	ms := make([]*bbolt2.Migration, 0, len(fs))
	for _, f := range fs {
		ms = append(ms, f())
	}
	m := bbolt2.NewMigrate(db, ms...)
	if err := m.Migrate(); err != nil {
		return err
	}
	return nil
}
