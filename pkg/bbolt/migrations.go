package bbolt

import (
	"errors"
	"fmt"
	"sort"

	"go.etcd.io/bbolt"
)

// MigrateFunc 数据库迁移函数
type MigrateFunc func(db *bbolt.DB) error

// RollbackFunc 数据库回滚函数
type RollbackFunc func(db *bbolt.DB) error

type Migration struct {
	ID       string
	Migrate  MigrateFunc
	Rollback RollbackFunc
}

type Migrate struct {
	db         *bbolt.DB
	migrations []*Migration
}

func NewMigrate(db *bbolt.DB, migrations ...*Migration) *Migrate {
	return &Migrate{
		db:         db,
		migrations: migrations,
	}
}

func (m *Migrate) Migrate() error {
	// 获取所有的迁移记录
	records := make([]string, 0)
	if err := m.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("migrations"))
		if err != nil {
			return err
		}
		return b.ForEach(func(k, v []byte) error {
			records = append(records, string(v))
			return nil
		})
	}); err != nil {
		return err
	}
	// 对比迁移记录和迁移函数
	if len(records) > len(m.migrations) {
		return errors.New("migrate records more than migrate functions")
	}
	// 排序
	sort.Strings(records)
	for n, record := range records {
		if record != m.migrations[n].ID {
			return fmt.Errorf("migrate id not match: %s != %s", record, m.migrations[n].ID)
		}
	}
	// 取出未迁移的函数
	migrations := m.migrations[len(records):]
	// 执行迁移
	for _, migration := range migrations {
		if err := migration.Migrate(m.db); err != nil {
			return err
		}
	}
	return nil
}
