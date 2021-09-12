package repository

import (
	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	"gorm.io/gorm"
)

type User interface {
	FindById(id int64) (*entity.User, error)
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	return &user{}
}

func (u *user) FindById(id int64) (*entity.User, error) {
	ret := &entity.User{ID: id}
	err := u.db.First(ret).Error
	return ret, err
}
