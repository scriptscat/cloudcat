package service

import (
	"strings"
	"time"

	"github.com/scriptscat/cloudcat/internal/domain/user/dto"
	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	"github.com/scriptscat/cloudcat/internal/domain/user/errs"
	"github.com/scriptscat/cloudcat/internal/domain/user/repository"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
	"github.com/scriptscat/cloudcat/pkg/utils"
)

type User interface {
	Login(login *dto.Login) (*dto.UserInfo, error)
	Register(register *dto.Register) (*dto.UserInfo, error)
	RequestRegisterEmailCode(email string) (*entity.VerifyCode, error)
	UserInfo(uid int64) (*dto.UserInfo, error)
	UploadAvatar(uid int64, b []byte) error
	oauthRegister(user *entity.User) (int64, error)
}

const (
	EnableRegister = "enable_register"
	EnableInvcode  = "enable_invcode"

	RequiredVerifyEmail = "required_verify_email"
	AllowEmailSuffix    = "allow_email_suffix"
)

type user struct {
	config     config.SystemConfig
	kv         kvdb.KvDb
	userRepo   repository.User
	verifyRepo repository.VerifyCode
}

func NewUser(config config.SystemConfig, kv kvdb.KvDb, userRepo repository.User, verifyRepo repository.VerifyCode) User {
	return &user{
		config:     config,
		userRepo:   userRepo,
		verifyRepo: verifyRepo,
		kv:         kv,
	}
}

func (u *user) UserInfo(uid int64) (*dto.UserInfo, error) {
	user, err := u.userRepo.FindById(uid)
	if err != nil {
		return nil, err
	}
	return u.toUserInfo(user)
}

func (u *user) toUserInfo(user *entity.User) (*dto.UserInfo, error) {
	if user == nil {
		return nil, errs.ErrUserNotFound
	}
	return dto.ToUserInfo(user), nil
}

func (u *user) Login(login *dto.Login) (*dto.UserInfo, error) {
	var user *entity.User
	var err error
	if strings.Index(login.Account, "@") == -1 {
		user, err = u.userRepo.FindByMobile(login.Account)
	} else {
		user, err = u.userRepo.FindByEmail(login.Account)
	}
	if err != nil {
		return nil, err
	}
	info, err := u.toUserInfo(user)
	if err != nil {
		return nil, err
	}
	if err := user.CheckPassword(login.Password); err != nil {
		return nil, err
	}
	return info, nil
}

func (u *user) Register(register *dto.Register) (*dto.UserInfo, error) {
	enable, err := u.config.GetConfig(EnableRegister)
	if err != nil {
		return nil, err
	}
	if enable == "0" {
		return nil, errs.ErrRegisterDisable
	}
	verifyEmail, err := u.config.GetConfig(RequiredVerifyEmail)
	if err != nil {
		return nil, err
	}
	if err := u.checkEmail(register.Email); err != nil {
		return nil, err
	}
	if verifyEmail == "1" {
		if register.EmailVerifyCode == "" {
			return nil, errs.ErrRegisterVerifyEmail
		}
		vcode, err := u.verifyRepo.FindById(register.Email)
		if err != nil {
			return nil, err
		}
		if err := vcode.CheckCode(register.EmailVerifyCode); err != nil {
			return nil, err
		}
	}
	user := &entity.User{
		Nickname:   register.Nickname,
		Email:      register.Email,
		Role:       "user",
		Createtime: time.Now().Unix(),
		Updatetime: 0,
	}
	if err := user.SetPassword(register.Password); err != nil {
		return nil, err
	}
	if err := u.userRepo.Save(user); err != nil {
		return nil, err
	}
	return dto.ToUserInfo(user), nil
}

func (u *user) checkMobile(mobile string) error {
	user, err := u.userRepo.FindByMobile(mobile)
	if err != nil {
		return err
	}
	if user != nil {
		return errs.ErrEmailExist
	}
	return nil
}

func (u *user) checkEmail(email string) error {
	emailSuffix, err := u.config.GetConfig(AllowEmailSuffix)
	if err != nil {
		return err
	}
	if emailSuffix != "" {
		suffixs := strings.Split(emailSuffix, ",")
		flag := false
		for _, v := range suffixs {
			if strings.HasSuffix(email, v) {
				flag = true
				break
			}
		}
		if !flag {
			return errs.ErrEmailSuffixNotAllow
		}
	}
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user != nil {
		return errs.ErrEmailExist
	}
	return nil
}

func (u *user) RequestRegisterEmailCode(email string) (*entity.VerifyCode, error) {
	if err := u.checkEmail(email); err != nil {
		return nil, err
	}
	v := &entity.VerifyCode{
		Identifier: email,
		Op:         "register",
		Code:       strings.ToUpper(utils.RandString(6, 0)),
		Expiretime: time.Now().Add(time.Minute * 5).Unix(),
	}
	if err := u.verifyRepo.Save(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (u *user) UploadAvatar(uid int64, b []byte) error {

	return nil
}

func (u *user) oauthRegister(user *entity.User) (int64, error) {
	if err := u.userRepo.Save(user); err != nil {
		return 0, nil
	}
	return user.ID, nil
}
