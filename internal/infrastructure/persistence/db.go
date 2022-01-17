package persistence

import (
	"github.com/scriptscat/cloudcat/internal/infrastructure/persistence/migrations"
	"github.com/scriptscat/cloudcat/internal/pkg/database"
	"github.com/scriptscat/cloudcat/internal/service/user/infrastructure/persistence"
)

type Repositories struct {
	db   *database.Database
	User *persistence.Repositories
}

func NewRepositories(db *database.Database) *Repositories {
	return &Repositories{
		db:   db,
		User: persistence.NewRepositories(db.DB),
	}
}

func (r *Repositories) Migrations() error {
	return migrations.RunMigrations(r.db)
}
