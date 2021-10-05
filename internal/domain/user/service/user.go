package service

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/scriptscat/cloudcat/internal/domain/user/dto"
	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	"github.com/scriptscat/cloudcat/internal/domain/user/errs"
	"github.com/scriptscat/cloudcat/internal/domain/user/repository"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

type User interface {
	Login(login *dto.Login) (*dto.UserInfo, error)
	Register(register *dto.Register) (*dto.UserInfo, error)
	RequestEmailCode(email, op string) (*entity.VerifyCode, error)
	UserInfo(uid int64) (*dto.UserInfo, error)
	Avatar(uid int64) ([]byte, error)
	UploadAvatar(uid int64, b []byte) error
	oauthRegister(tx *gorm.DB, user *entity.User) (int64, error)
	CheckUsername(username string) error
	UpdateUserInfo(uid int64, req *dto.UpdateUserInfo) error
	UpdatePassword(uid int64, req *dto.UpdatePassword) error
}

const (
	EnableRegister = "enable_register"
	EnableInvcode  = "enable_invcode"

	RequiredVerifyEmail = "required_verify_email"
	AllowEmailSuffix    = "allow_email_suffix"
)

type user struct {
	config      config.SystemConfig
	kv          kvdb.KvDb
	userRepo    repository.User
	verifyRepo  repository.VerifyCode
	resourceDir string
}

func NewUser(config config.SystemConfig, kv kvdb.KvDb, userRepo repository.User, verifyRepo repository.VerifyCode) User {
	return &user{
		config:      config,
		userRepo:    userRepo,
		verifyRepo:  verifyRepo,
		kv:          kv,
		resourceDir: "./resource/user",
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
	if strings.Index(login.Username, "@") != -1 && strings.Index(login.Username, ".") != -1 {
		user, err = u.userRepo.FindByEmail(login.Username)
	} else {
		user, err = u.userRepo.FindByName(login.Username)
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
	if err := u.CheckUsername(register.Username); err != nil {
		return nil, err
	}
	if err := u.checkEmail(register.Email); err != nil {
		return nil, err
	}
	verifyEmail, err := u.config.GetConfig(RequiredVerifyEmail)
	if err != nil {
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
		if vcode == nil {
			return nil, errs.ErrEmailVCodeNotFound
		}
		if err := vcode.CheckCode(register.EmailVerifyCode, "register"); err != nil {
			return nil, err
		}
	}
	user := &entity.User{
		Username:   register.Username,
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
		return errs.ErrMobileExist
	}
	return nil
}

func (u *user) CheckUsername(username string) error {
	user, err := u.userRepo.FindByName(username)
	if err != nil {
		return err
	}
	if user != nil {
		return errs.ErrUsernameExist
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

func (u *user) RequestEmailCode(email, op string) (*entity.VerifyCode, error) {
	if err := u.checkEmail(email); err != nil {
		return nil, err
	}
	v := &entity.VerifyCode{
		Identifier: email,
		Op:         op,
		Code:       strings.ToUpper(utils.RandString(6, 0)),
		Expiretime: time.Now().Add(time.Minute * 5).Unix(),
	}
	if err := u.verifyRepo.Save(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (u *user) Avatar(uid int64) ([]byte, error) {
	user, err := u.userRepo.FindById(uid)
	if err != nil {
		return nil, err
	}
	if user.Avatar == "" {
		return nil, errs.ErrAvatarIsNil
	}
	return os.ReadFile(user.Avatar)
}

// UploadAvatar NOTE: resource这块可能要重构
func (u *user) UploadAvatar(uid int64, b []byte) error {
	user, err := u.UserInfo(uid)
	if err != nil {
		return err
	}
	c := http.DetectContentType(b)
	if strings.Index(c, "image") == -1 {
		return errs.ErrAvatarNotImage
	}
	if len(b) > 1024*1024 {
		return errs.ErrAvatarTooBig
	}
	suffix := c[strings.LastIndex(c, "/")+1:]
	p, name := u.getDir(b, "."+suffix)
	p = path.Join(u.resourceDir, "avatar", p)
	if err := os.MkdirAll(p, 0755); err != nil {
		return err
	}
	p = path.Join(p, name)
	if err := os.WriteFile(p, b, 0644); err != nil {
		return err
	}
	return u.userRepo.SaveUserAvatar(user.ID, p)
}

func (u *user) getDir(b []byte, suffix string) (string, string) {
	d := fmt.Sprintf("%x", sha1.Sum(b))
	return path.Join(d[:2], d[2:4]), d + suffix
}

//NOTE: 有点丑陋,先简单实现了
func (u *user) oauthRegister(tx *gorm.DB, user *entity.User) (int64, error) {
	userRepo := repository.NewUser(tx)
	if err := userRepo.Save(user); err != nil {
		return 0, nil
	}
	return user.ID, nil
}

func (u *user) UpdateUserInfo(uid int64, req *dto.UpdateUserInfo) error {
	user, err := u.userRepo.FindById(uid)
	if err != nil {
		return err
	}
	if req.Username != user.Username {
		// 更新用户名
		if u, err := u.userRepo.FindByName(req.Username); err != nil {
			return err
		} else if u != nil {
			return errs.ErrUsernameExist
		}
		user.Username = req.Username
	}
	if req.Email != req.Email {
		// 更新邮箱
		vcode, err := u.verifyRepo.FindById(req.Email)
		if err != nil {
			return err
		}
		if vcode == nil {
			return errs.ErrEmailVCodeNotFound
		}
		if err := vcode.CheckCode(req.EmailVerifyCode, "change-user-info"); err != nil {
			return err
		}
		user.Email = req.Email
	}
	return u.userRepo.Save(user)
}

func (u *user) UpdatePassword(uid int64, req *dto.UpdatePassword) error {
	user, err := u.userRepo.FindById(uid)
	if err != nil {
		return err
	}
	if err := user.CheckPassword(req.Password); err != nil {
		return err
	}
	if err := user.SetPassword(req.Password); err != nil {
		return err
	}
	return u.userRepo.Save(user)
}
