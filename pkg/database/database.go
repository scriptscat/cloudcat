package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Database struct {
	*gorm.DB
	cfg *Config
}

func NewDatabase(cfg *Config) (*Database, error) {
	var err error
	db := &Database{cfg: cfg}
	switch cfg.Type {
	case "mysql":
		db.DB, err = gorm.Open(mysql.New(mysql.Config{
			DSN: cfg.Mysql.Dsn,
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   cfg.Mysql.Prefix,
				SingularTable: true,
			},
		})
	case "sqlite":
		db.DB, err = gorm.Open(sqlite.Open(cfg.Sqlite.File), &gorm.Config{})
	}
	if err != nil {
		return nil, err
	}
	return db, nil
}
