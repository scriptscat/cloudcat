package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
)

type VerifyCode interface {
	Save(vcode *entity.VerifyCode) error
	FindById(id string) (*entity.VerifyCode, error)
}

type verifyCode struct {
	kv kvdb.KvDb
}

func NewVerifyCode(kv kvdb.KvDb) VerifyCode {
	return &verifyCode{kv: kv}
}

func (v *verifyCode) Save(vcode *entity.VerifyCode) error {
	j, err := json.Marshal(vcode)
	if err != nil {
		return err
	}
	return v.kv.Set(context.Background(), v.key(vcode.Identifier), string(j), time.Hour*24)
}

func (v *verifyCode) FindById(id string) (*entity.VerifyCode, error) {
	s, err := v.kv.Get(context.Background(), v.key(id))
	if err != nil {
		return nil, err
	}
	ret := &entity.VerifyCode{}
	if err := json.Unmarshal([]byte(s), ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (v *verifyCode) key(id string) string {
	return "user:verify:code:" + id
}
