package service

import (
	"strings"
	"sync"
	"time"

	"github.com/scriptscat/cloudcat/internal/domain/user/dto"
	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	"github.com/scriptscat/cloudcat/internal/domain/user/errs"
	"github.com/scriptscat/cloudcat/internal/domain/user/repository"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
	"github.com/scriptscat/cloudcat/pkg/oauth/bbs"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/silenceper/wechat/v2"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/silenceper/wechat/v2/util"
)

type User interface {
	Login(login *dto.Login) (*dto.UserInfo, error)
	Register(register *dto.Register) (*dto.UserInfo, error)
	RequestRegisterEmailCode(email string) (*entity.VerifyCode, error)

	RedirectOAuth(redirectUrl, platform string) (string, error)
	BBSOAuthLogin(code string) (*dto.OAuthRespond, error)
	WechatAuthLogin(code string) (*dto.OAuthRespond, error)
	UserInfo(uid int64) (*dto.UserInfo, error)
}

const (
	ENABLE_REGISTER = "enable_register"
	ENABLE_INVCODE  = "enable_invcode"

	OAUTH_CONFIG_BBS_CLIENT_ID     = "oauth_config_bbs_client_id"
	OAUTH_CONFIG_BBS_CLIENT_SECRET = "oauth_config_bbs_client_secret"

	OAUTH_CONFIG_WECHAT_APP_ID         = "oauth_config_wechat_app_id"
	OAUTH_CONFIG_WECHAT_APP_SECRET     = "oauth_config_wechat_app_secret"
	OAUTH_CONFIG_WECHAT_TOKEN          = "oauth_config_wechat_token"
	OAUTH_CONFIG_WECHAT_ENCODINGAESKEY = "oauth_config_wechat_encoding_aes_key"

	REQUIRED_VERIFY_EMAIL = "required_verify_email"
	ALLOW_EMAIL_SUFFIX    = "allow_email_suffix"
)

type user struct {
	kv          kvdb.KvDb
	config      config.SystemConfig
	userRepo    repository.User
	verifyRepo  repository.VerifyCode
	bbsOAuth    repository.BBSOAuth
	wechatOAuth repository.WechatOAuth
	oauth       struct {
		sync.RWMutex
		bbs    *bbs.Client
		wechat *oauth.Oauth
	}
}

func NewUser(config config.SystemConfig, kv kvdb.KvDb, bbs repository.BBSOAuth, wc repository.WechatOAuth, userRepo repository.User, verifyRepo repository.VerifyCode) User {
	return &user{
		config:      config,
		bbsOAuth:    bbs,
		wechatOAuth: wc,
		userRepo:    userRepo,
		verifyRepo:  verifyRepo,
		kv:          kv,
	}
}

func (u *user) RedirectOAuth(redirectUrl, platform string) (string, error) {
	switch platform {
	case "bbs":
		client, err := u.getBbsClient()
		if err != nil {
			return "", err
		}
		return client.RedirectURL(redirectUrl), err
	case "wechat":
		client, err := u.getWechatClient()
		if err != nil {
			return "", err
		}
		return client.GetRedirectURL(redirectUrl, "snsapi_userinfo", "")
	}
	return "", errs.ErrOAuthPlatformNotSupport
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

func (u *user) BBSOAuthLogin(code string) (*dto.OAuthRespond, error) {
	client, err := u.getBbsClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.RequestAccessToken(code)
	if err != nil {
		return nil, err
	}
	userResp, err := client.RequestUser(resp.AccessToken)
	if err != nil {
		return nil, err
	}
	bbs, err := u.bbsOAuth.FindByOpenid(userResp.User.Uid)
	if err != nil {
		return nil, err
	}
	if bbs == nil {
		// 需要绑定账号登录
		return &dto.OAuthRespond{
			IsBind: false,
		}, nil
	}
	user, err := u.UserInfo(bbs.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.OAuthRespond{
		UserInfo: user,
		IsBind:   true,
	}, nil
}

func (u *user) WechatAuthLogin(code string) (*dto.OAuthRespond, error) {
	client, err := u.getWechatClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetUserAccessToken(code)
	if err != nil {
		return nil, err
	}
	wechat, err := u.wechatOAuth.FindByOpenid(resp.OpenID)
	if err != nil {
		return nil, err
	}
	if wechat == nil {
		// 需要绑定账号登录
		return &dto.OAuthRespond{
			IsBind: false,
		}, nil
	}
	user, err := u.UserInfo(wechat.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.OAuthRespond{
		UserInfo: user,
		IsBind:   true,
	}, nil
}

func (u *user) getWechatClient() (*oauth.Oauth, error) {
	u.oauth.RLock()
	if u.oauth.wechat == nil {
		u.oauth.RUnlock()
		u.oauth.Lock()
		defer u.oauth.Unlock()
		wc := wechat.NewWechat()
		appId, err := u.getOAuthConfig(OAUTH_CONFIG_WECHAT_APP_ID)
		if err != nil {
			return nil, err
		}
		appSecret, err := u.getOAuthConfig(OAUTH_CONFIG_WECHAT_APP_SECRET)
		if err != nil {
			return nil, err
		}
		token, err := u.getOAuthConfig(OAUTH_CONFIG_WECHAT_APP_SECRET)
		if err != nil {
			return nil, err
		}
		encodingAESKey, err := u.getOAuthConfig(OAUTH_CONFIG_WECHAT_APP_SECRET)
		if err != nil {
			return nil, err
		}
		sysCache := utils.NewWxCache(u.kv)
		of := wc.GetOfficialAccount(&offConfig.Config{
			AppID:          appId,
			AppSecret:      appSecret,
			Token:          token,
			EncodingAESKey: encodingAESKey,
			Cache:          sysCache,
		})
		u.oauth.wechat = of.GetOauth()
	} else {
		u.oauth.RUnlock()
	}
	return u.oauth.wechat, nil

}

func (u *user) getBbsClient() (*bbs.Client, error) {
	u.oauth.RLock()
	if u.oauth.wechat == nil {
		u.oauth.RUnlock()
		u.oauth.Lock()
		defer u.oauth.Unlock()
		clientid, err := u.getOAuthConfig(OAUTH_CONFIG_BBS_CLIENT_ID)
		if err != nil {
			return nil, err
		}
		clientSecret, err := u.getOAuthConfig(OAUTH_CONFIG_BBS_CLIENT_SECRET)
		if err != nil {
			return nil, err
		}
		u.oauth.bbs = bbs.NewClient(bbs.Config{
			ClientID:     clientid,
			ClientSecret: clientSecret,
		})
	} else {
		u.oauth.RUnlock()
	}
	return u.oauth.bbs, nil
}

func (u *user) getOAuthConfig(key string) (string, error) {
	ret, err := u.config.GetConfig(key)
	if err != nil {
		return "", err
	}
	if ret == "" {
		return "", errs.ErrOAuthPlatformNotConfigured
	}
	return ret, nil
}

func (u *user) Login(login *dto.Login) (*dto.UserInfo, error) {
	var user *entity.User
	var err error
	if login.Username != "" {
		user, err = u.userRepo.FindByName(login.Username)
	} else if login.Email != "" {
		user, err = u.userRepo.FindByEmail(login.Email)
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
	enable, err := u.config.GetConfig(ENABLE_REGISTER)
	if err != nil {
		return nil, err
	}
	if enable == "0" {
		return nil, errs.ErrRegisterDisable
	}
	verifyEmail, err := u.config.GetConfig(REQUIRED_VERIFY_EMAIL)
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

func (u *user) checkEmail(email string) error {
	emailSuffix, err := u.config.GetConfig(ALLOW_EMAIL_SUFFIX)
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
		Code:       strings.ToUpper(util.RandomStr(6)),
		Expiretime: time.Now().Add(time.Minute * 5).Unix(),
	}
	if err := u.verifyRepo.Save(v); err != nil {
		return nil, err
	}
	return v, nil
}
