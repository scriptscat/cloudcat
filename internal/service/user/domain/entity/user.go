package entity

import (
	"github.com/scriptscat/cloudcat/internal/service/user/domain/errs"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/vo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64  `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"unique;column:username;type:varchar(128);not null" json:"username"` // 用户名
	PasswordHash string `gorm:"column:password_hash;type:varchar(200)" json:"password_hash"`
	Email        string `gorm:"unique;column:email;type:varchar(128);default:null" json:"email"`
	Mobile       string `gorm:"unique;column:mobile;type:varchar(128);default:null" json:"mobile"`
	Avatar       string `gorm:"column:avatar;type:varchar(128)" json:"avatar"`
	Role         string `gorm:"column:role;type:varchar(16);not null" json:"role"`
	Createtime   int64  `gorm:"column:createtime;type:bigint(20);not null" json:"createtime"`
	Updatetime   int64  `gorm:"column:updatetime;type:bigint(20)" json:"updatetime"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return errs.ErrWrongPassword
	}
	return nil
}

func (u *User) PublicUser() *vo.UserInfo {
	return &vo.UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
		Role:     u.Role,
	}
}
