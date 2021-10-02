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

func NewDatabase(cfg *Config, debug bool) (*Database, error) {
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
		if err == nil {
			db.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
		}
	case "sqlite":
		db.DB, err = gorm.Open(sqlite.Open(cfg.Sqlite.File), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	}
	if err != nil {
		return nil, err
	}
	if debug {
		db.DB = db.DB.Debug()
	}
	return db, nil
}
