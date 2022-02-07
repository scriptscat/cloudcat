package application

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	config2 "github.com/scriptscat/cloudcat/internal/infrastructure/config"
	"github.com/scriptscat/cloudcat/internal/infrastructure/sender"
	"github.com/scriptscat/cloudcat/internal/pkg/kvdb"
	entity2 "github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/errs"
	repository2 "github.com/scriptscat/cloudcat/internal/service/user/domain/repository"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/vo"
	persistence2 "github.com/scriptscat/cloudcat/internal/service/user/infrastructure/persistence"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"gorm.io/gorm"
)

type User interface {
	Login(login *vo.Login) (*vo.UserInfo, error)
	Register(register *vo.Register) (*vo.UserInfo, error)
	RequestEmailCode(email, op string) (*entity2.VerifyCode, error)
	UserInfo(uid int64) (*vo.UserInfo, error)
	Avatar(uid int64) ([]byte, error)
	UploadAvatar(uid int64, b []byte) error
	oauthRegister(tx *gorm.DB, user *entity2.User) (int64, error)
	CheckUsername(username string) error
	UpdateUserInfo(uid int64, req *vo.UpdateUserInfo) error
	UpdatePassword(uid int64, req *vo.UpdatePassword) error
	RequestForgetPasswordEmail(email string) error
	ValidResetPassword(code string) (*vo.UserInfo, error)
	ResetPassword(code, password string) error
	UpdateEmail(uid int64, req *vo.UpdateEmail) error
}

const (
	EnableRegister = "enable_register"
	EnableInvcode  = "enable_invcode"

	RequiredVerifyEmail = "required_verify_email"
	AllowEmailSuffix    = "allow_email_suffix"
)

type user struct {
	config      config2.SystemConfig
	kv          kvdb.KvDb
	userRepo    repository2.User
	verifyRepo  repository2.VerifyCode
	resourceDir string
	sender      sender.Sender
}

func NewUser(config config2.SystemConfig, kv kvdb.KvDb, userRepo repository2.User, verifyRepo repository2.VerifyCode, sender sender.Sender) User {
	return &user{
		config:      config,
		userRepo:    userRepo,
		verifyRepo:  verifyRepo,
		kv:          kv,
		resourceDir: "./resource/user",
		sender:      sender,
	}
}

func (u *user) UserInfo(uid int64) (*vo.UserInfo, error) {
	user, err := u.userRepo.FindById(uid)
	if err != nil {
		return nil, err
	}
	return user.PublicUser(), nil
}

func (u *user) Login(login *vo.Login) (*vo.UserInfo, error) {
	var user *entity2.User
	var err error
	if strings.Index(login.Username, "@") != -1 && strings.Index(login.Username, ".") != -1 {
		user, err = u.userRepo.FindByEmail(login.Username)
	} else {
		user, err = u.userRepo.FindByName(login.Username)
	}
	if err != nil {
		return nil, err
	}
	info := user.PublicUser()
	if err != nil {
		return nil, err
	}
	if err := user.CheckPassword(login.Password); err != nil {
		return nil, err
	}
	return info, nil
}

func (u *user) Register(register *vo.Register) (*vo.UserInfo, error) {
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
		if err := u.verifyRepo.InvalidCode(vcode); err != nil {
			return nil, err
		}
	}
	user := &entity2.User{
		Username:   register.Username,
		Email:      register.Email,
		Role:       "user",
		Createtime: time.Now().Unix(),
		Updatetime: 0,
	}
	if err := user.Register(register.Password); err != nil {
		return nil, err
	}
	if err := u.userRepo.SaveUser(user); err != nil {
		return nil, err
	}
	return user.PublicUser(), nil
}

func (u *user) checkMobile(mobile string) error {
	_, err := u.userRepo.FindByMobile(mobile)
	if err != nil {
		if err == errs.ErrUserNotFound {
			return nil
		}
		return err
	}
	return errs.ErrMobileExist
}

func (u *user) CheckUsername(username string) error {
	_, err := u.userRepo.FindByName(username)
	if err != nil {
		if err == errs.ErrUserNotFound {
			return nil
		}
		return err
	}
	return errs.ErrUsernameExist
}

func (u *user) checkEmail(email string) error {
	emailSuffix, err := u.config.GetConfig(AllowEmailSuffix)
	if err != nil {
		return err
	}
	if emailSuffix != "" {
		suffixs := strings.Split(emailSuffix, ",")
		flag := false
		s := strings.Split(email, "@")
		if len(s) != 2 {
			return errs.ErrEmailExist
		}
		for _, v := range suffixs {
			if s[1] == v {
				flag = true
				break
			}
		}
		if !flag {
			return errs.ErrEmailSuffixNotAllow
		}
	}
	_, err = u.userRepo.FindByEmail(email)
	if err != nil {
		if err == errs.ErrUserNotFound {
			return nil
		}
		return err
	}
	return errs.ErrEmailExist
}

func (u *user) RequestEmailCode(email, op string) (*entity2.VerifyCode, error) {
	if err := u.checkEmail(email); err != nil {
		return nil, err
	}
	v := &entity2.VerifyCode{
		Identifier: email,
		Op:         op,
		Code:       strings.ToUpper(utils.RandString(6, 0)),
		Expired:    time.Now().Add(time.Minute * 5).Unix(),
	}
	if err := u.verifyRepo.SaveVerifyCode(v); err != nil {
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
func (u *user) oauthRegister(tx *gorm.DB, user *entity2.User) (int64, error) {
	userRepo := persistence2.NewUser(tx)
	if err := userRepo.SaveUser(user); err != nil {
		return 0, nil
	}
	return user.ID, nil
}

func (u *user) UpdateUserInfo(uid int64, req *vo.UpdateUserInfo) error {
	user, err := u.userRepo.FindById(uid)
	if err != nil {
		return err
	}
	if err := u.CheckUsername(req.Username); err != nil {
		return err
	}
	user.Username = req.Username
	return u.userRepo.SaveUser(user)
}

func (u *user) UpdatePassword(uid int64, req *vo.UpdatePassword) error {
	user, err := u.userRepo.FindById(uid)
	if err != nil {
		return err
	}
	if err := user.UpdatePassword(req.OldPassword, req.Password); err != nil {
		return err
	}
	return u.userRepo.SaveUser(user)
}

func (u *user) RequestForgetPasswordEmail(email string) error {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	vcode := &entity2.VerifyCode{
		Identifier: user.Email,
		Op:         "forget-password",
		Code:       utils.RandString(32, 2),
		Expired:    time.Now().Add(time.Minute * 30).Unix(),
	}
	if err := u.verifyRepo.SaveVerifyCode(vcode); err != nil {
		return err
	}
	url, err := u.config.GetConfig(config2.HomeUrl)
	if err != nil {
		return err
	}
	url += "/user/reset-password?code=" + vcode.Code
	return u.sender.SendEmail(user.Email, "找回密码", "请点击链接<a href=\""+url+"\">找回密码</a>或者复制链接: "+url+" 进行访问 链接有效期30分钟,请在30分钟内使用", "text/html")
}

func (u *user) ValidResetPassword(code string) (*vo.UserInfo, error) {
	vcode, err := u.verifyRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}
	if err := vcode.CheckCode(code, "forget-password"); err != nil {
		return nil, err
	}
	user, err := u.userRepo.FindByEmail(vcode.Identifier)
	if err != nil {
		return nil, err
	}
	return user.PublicUser(), nil
}

func (u *user) ResetPassword(code, password string) error {
	vcode, err := u.verifyRepo.FindByCode(code)
	if err != nil {
		return err
	}
	user, err := u.userRepo.FindByEmail(vcode.Identifier)
	if err != nil {
		return err
	}
	if err := user.ResetPassword(vcode, code, password); err != nil {
		return err
	}
	if err := u.verifyRepo.InvalidCode(vcode); err != nil {
		return err
	}
	return u.userRepo.SaveUser(user)
}

func (u *user) UpdateEmail(uid int64, req *vo.UpdateEmail) error {
	user, err := u.userRepo.FindById(uid)
	if err != nil {
		return err
	}
	if req.Email != user.Email {
		if err := u.checkEmail(req.Email); err != nil {
			return err
		}
		vcode, err := u.verifyRepo.FindById(req.Email)
		if err != nil {
			return err
		}
		// 更新邮箱
		if err := user.UpdateEmail(vcode, req.Code, req.Email); err != nil {
			return err
		}
		if err := u.verifyRepo.InvalidCode(vcode); err != nil {
			return err
		}
	} else {
		return errs.ErrEmailExist
	}
	return u.userRepo.SaveUser(user)
}
