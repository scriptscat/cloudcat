package persistence

import (
	"net/http"

	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/repository"
	"github.com/scriptscat/cloudcat/pkg/errs"
	"gorm.io/gorm"
)

type verifyCode struct {
	db *gorm.DB
}

func NewVerifyCode(db *gorm.DB) repository.VerifyCode {
	return &verifyCode{
		db: db,
	}
}

func (v *verifyCode) SaveVerifyCode(vcode *entity.VerifyCode) error {
	// 删除原来的
	if err := v.db.Delete(vcode, "identifier=?", vcode.Identifier).Error; err != nil {
		return err
	}
	return v.db.Save(vcode).Error
}

func (v *verifyCode) FindById(id string) (*entity.VerifyCode, error) {
	ret := &entity.VerifyCode{}
	if err := v.db.First(ret, "identifier=?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewError(http.StatusOK, 1000, "验证码错误")
		}
		return nil, err
	}
	return ret, nil
}

func (v *verifyCode) FindByCode(code string) (*entity.VerifyCode, error) {
	ret := &entity.VerifyCode{}
	if err := v.db.First(ret, "code=?", code).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewError(http.StatusOK, 1001, "验证码错误")
		}
		return nil, err
	}
	return ret, nil
}

func (v *verifyCode) InvalidCode(vcode *entity.VerifyCode) error {
	return v.db.Delete(vcode, "id=?", vcode.ID).Error
}
