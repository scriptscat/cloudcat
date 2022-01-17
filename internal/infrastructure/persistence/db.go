package persistence

import (
	"github.com/scriptscat/cloudcat/internal/service/user/infrastructure/persistence"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

type Repositories struct {
	User *persistence.Repositories
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User: persistence.NewRepositories(db),
	}
}

func (r *Repositories) AutoMigrate() error {
	return utils.Errs(
		r.User.AutoMigrate(),
	)
}
