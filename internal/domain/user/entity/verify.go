package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/scriptscat/cloudcat/internal/pkg/errs"
)

type VerifyCode struct {
	Identifier string
	Op         string
	Code       string
	Expiretime int64
}

func (v *VerifyCode) CheckCode(code string) error {
	if time.Now().Unix() > v.Expiretime {
		return errs.NewBadRequestError(1000, "验证码过期")
	}
	if v.Code == "" {
		return errors.New("验证码怎么会为空")
	}
	if v.Code != strings.ToUpper(code) {
		return errs.NewBadRequestError(1001, "验证码错误")
	}
	return nil
}
