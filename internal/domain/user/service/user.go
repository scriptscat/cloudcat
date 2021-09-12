package service

import (
	"sync"

	"github.com/scriptscat/cloudcat/internal/domain/user/dto"
	"github.com/scriptscat/cloudcat/internal/domain/user/errs"
	"github.com/scriptscat/cloudcat/internal/domain/user/repository"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
	"github.com/scriptscat/cloudcat/pkg/oauth/bbs"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/silenceper/wechat/v2"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
)

type User interface {
	RedirectOAuth(redirectUrl, platform string) (string, error)
	BBSOAuthLogin(code string) (*dto.OAuthRespond, error)
	WechatAuthLogin(code string) (*dto.OAuthRespond, error)
}

const (
	OAUTH_CONFIG_BBS_CLIENT_ID     = "oauth_config_bbs_client_id"
	OAUTH_CONFIG_BBS_CLIENT_SECRET = "oauth_config_bbs_client_secret"

	OAUTH_CONFIG_WECHAT_APP_ID         = "oauth_config_wechat_app_id"
	OAUTH_CONFIG_WECHAT_APP_SECRET     = "oauth_config_wechat_app_secret"
	OAUTH_CONFIG_WECHAT_TOKEN          = "oauth_config_wechat_token"
	OAUTH_CONFIG_WECHAT_ENCODINGAESKEY = "oauth_config_wechat_encoding_aes_key"

	REQUIRED_EMAIL        = "REQUIRED_EMAIL"
	REQUIRED_VERIFY_EMAIL = "REQUIRED_VERIFY_EMAIL"
)

type user struct {
	kv          kvdb.KvDb
	config      config.SystemConfig
	userRepo    repository.User
	bbsOAuth    repository.BBSOAuth
	wechatOAuth repository.WechatOAuth
	oauth       struct {
		sync.RWMutex
		bbs    *bbs.Client
		wechat *oauth.Oauth
	}
}

func NewUser(config config.SystemConfig, kv kvdb.KvDb, bbs repository.BBSOAuth, wc repository.WechatOAuth, userRepo repository.User) User {
	return &user{
		config:      config,
		bbsOAuth:    bbs,
		wechatOAuth: wc,
		userRepo:    userRepo,
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
