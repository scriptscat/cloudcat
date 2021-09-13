package entity

import (
	"github.com/scriptscat/cloudcat/internal/domain/user/errs"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64  `gorm:"primaryKey" json:"id"`
	Nickname     string `gorm:"index:nickname;column:nickname;type:varchar(255);not null" json:"nickname"` // 用户名
	PasswordHash string `gorm:"column:password_hash;type:varchar(255)" json:"password_hash"`
	Email        string `gorm:"unique;column:email;type:varchar(255)" json:"email"`
	Mobile       string `gorm:"unique;column:mobile;type:varchar(255)" json:"mobile"`
	Avatar       string `gorm:"column:avatar;type:varchar(255)" json:"avatar"`
	Role         string `gorm:"column:role;type:varchar(255);not null" json:"role"`
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
