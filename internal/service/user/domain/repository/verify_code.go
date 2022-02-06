package repository

import (
	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
)

//go:generate mockgen -source ./verify_code.go -destination ./mock/verify_code.go
type VerifyCode interface {
	SaveVerifyCode(vcode *entity.VerifyCode) error
	FindById(id string) (*entity.VerifyCode, error)
	FindByCode(code string) (*entity.VerifyCode, error)
	InvalidCode(vcode *entity.VerifyCode) error
}
