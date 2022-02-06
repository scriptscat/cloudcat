package repository

import (
	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
)

//go:generate mockgen -source ./user.go -destination ./mock/user.go
type User interface {
	SaveUser(user *entity.User) error
	SaveUserAvatar(id int64, avatar string) error
	FindById(id int64) (*entity.User, error)
	FindByName(name string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByMobile(mobile string) (*entity.User, error)
}
