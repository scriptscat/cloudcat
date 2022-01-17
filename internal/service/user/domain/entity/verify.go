package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/scriptscat/cloudcat/pkg/errs"
)

type VerifyCode struct {
	ID         int64  `gorm:"primaryKey" json:"id"`
	Identifier string `gorm:"column:identifier;type:varchar(255);index;NOT NULL" json:"identifier"`
	Op         string `gorm:"column:op;type:varchar(255);NOT NULL" json:"op"`
	Code       string `gorm:"column:code;type:varchar(255);uniqueIndex;NOT NULL" json:"code"`
	Expired    int64  `gorm:"column:expired;type:bigint(255);NOT NULL" json:"expired"`
}

func (v *VerifyCode) CheckCode(code, op string) error {
	if time.Now().Unix() > v.Expired {
		return errs.NewBadRequestError(1000, "验证码过期")
	}
	if v.Code == "" {
		return errors.New("验证码怎么会为空")
	}
	if v.Code != strings.ToUpper(code) {
		return errs.NewBadRequestError(1001, "验证码错误")
	}
	if op != v.Op {
		return errs.NewBadRequestError(1002, "验证码错误")
	}
	return nil
}
