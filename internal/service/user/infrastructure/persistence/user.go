package persistence

import (
	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/repository"
	"gorm.io/gorm"
)

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) repository.User {
	return &user{
		db: db,
	}
}

func (u *user) Save(user *entity.User) error {
	return u.db.Save(user).Error
}

func (u *user) SaveUserAvatar(id int64, avatar string) error {
	return u.db.Model(&entity.User{}).Where("id=?", id).Update("avatar", avatar).Error
}

func (u *user) FindById(id int64) (*entity.User, error) {
	ret := &entity.User{ID: id}
	err := u.db.First(ret).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, err
}

func (u *user) FindByName(name string) (*entity.User, error) {
	ret := &entity.User{}
	err := u.db.Where("username=?", name).First(ret).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, err
}

func (u *user) FindByEmail(email string) (*entity.User, error) {
	ret := &entity.User{}
	err := u.db.Where("email=?", email).First(ret).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, err
}

func (u *user) FindByMobile(mobile string) (*entity.User, error) {
	ret := &entity.User{}
	err := u.db.Where("mobile=?", mobile).First(ret).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, err
}
