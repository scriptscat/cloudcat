package service

import (
	"errors"
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
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	offOAuth "github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/silenceper/wechat/v2/util"
	"github.com/sirupsen/logrus"
)

type OAuth interface {
	RedirectOAuth(redirectUrl, platform string) (string, error)
	BBSOAuthLogin(code string) (*dto.OAuthRespond, error)
	WechatAuthLogin(code string) (*dto.OAuthRespond, error)

	WechatScanLoginRequest() (*dto.WechatScanLogin, error)
	WechatScanLogin(openid, code string) error
	WechatScanLoginStatus(code string) (*dto.OAuthRespond, error)

	GetWechat() (*WechatConfig, error)
}

const (
	OAuthConfigBbsClientId     = "oauth_config_bbs_client_id"
	OAuthConfigBbsClientSecret = "oauth_config_bbs_client_secret"

	OAuthConfigWechatAppId          = "oauth_config_wechat_app_id"
	OAuthConfigWechatAppSecret      = "oauth_config_wechat_app_secret"
	OAuthConfigWechatToken          = "oauth_config_wechat_token"
	OAuthConfigWechatEncodingaeskey = "oauth_config_wechat_encoding_aes_key"
	OAuthConfigWechatReverseProxy   = "oauth_config_wechat_reverse_proxy"
)

type WechatConfig struct {
	Officialaccount *officialaccount.OfficialAccount
	ReverseProxy    string
}

type oauth struct {
	sync.RWMutex
	userSvc         User
	kv              kvdb.KvDb
	config          config.SystemConfig
	bbsOAuthRepo    repository.BBSOAuth
	wechatOAuthRepo repository.WechatOAuth

	wechat          *offOAuth.Oauth
	bbs             *bbs.Client
	officialaccount *officialaccount.OfficialAccount
}

func NewOAuth(config config.SystemConfig, kv kvdb.KvDb, userSvc User, bbs repository.BBSOAuth, wc repository.WechatOAuth) OAuth {
	return &oauth{
		config:          config,
		bbsOAuthRepo:    bbs,
		wechatOAuthRepo: wc,
		kv:              kv,
		userSvc:         userSvc,
	}
}

func (o *oauth) RedirectOAuth(redirectUrl, platform string) (string, error) {
	switch platform {
	case "bbs":
		client, err := o.getBbsClient()
		if err != nil {
			return "", err
		}
		return client.RedirectURL(redirectUrl), err
	case "wechat":
		client, err := o.getWechatClient()
		if err != nil {
			return "", err
		}
		return client.GetOauth().GetRedirectURL(redirectUrl, "snsapi_userinfo", "")
	}
	return "", errs.ErrOAuthPlatformNotSupport
}

func (o *oauth) BBSOAuthLogin(code string) (*dto.OAuthRespond, error) {
	client, err := o.getBbsClient()
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
	bbs, err := o.bbsOAuthRepo.FindByOpenid(userResp.User.Uid)
	if err != nil {
		return nil, err
	}
	if bbs == nil {
		// 需要绑定账号登录
		return &dto.OAuthRespond{
			IsBind: false,
		}, nil
	}
	user, err := o.userSvc.UserInfo(bbs.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.OAuthRespond{
		UserInfo: user,
		IsBind:   true,
	}, nil
}

func (o *oauth) WechatAuthLogin(code string) (*dto.OAuthRespond, error) {
	client, err := o.getWechatClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetOauth().GetUserAccessToken(code)
	if err != nil {
		return nil, err
	}
	wechat, err := o.wechatOAuthRepo.FindByOpenid(resp.OpenID)
	if err != nil {
		return nil, err
	}
	if wechat == nil {
		// 需要绑定账号登录
		return &dto.OAuthRespond{
			IsBind: false,
		}, nil
	}
	user, err := o.userSvc.UserInfo(wechat.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.OAuthRespond{
		UserInfo: user,
		IsBind:   true,
	}, nil
}

func (o *oauth) getWechatClient() (*officialaccount.OfficialAccount, error) {
	//TODO: 配置更新通知变更
	o.RLock()
	if o.wechat == nil {
		o.RUnlock()
		o.Lock()
		defer o.Unlock()
		wc := wechat.NewWechat()
		appId, err := o.getOAuthConfig(OAuthConfigWechatAppId)
		if err != nil {
			return nil, err
		}
		appSecret, err := o.getOAuthConfig(OAuthConfigWechatAppSecret)
		if err != nil {
			return nil, err
		}
		token, err := o.getOAuthConfig(OAuthConfigWechatAppSecret)
		if err != nil {
			return nil, err
		}
		encodingAESKey, err := o.getOAuthConfig(OAuthConfigWechatAppSecret)
		if err != nil {
			return nil, err
		}
		sysCache := utils.NewWxCache(o.kv)
		of := wc.GetOfficialAccount(&offConfig.Config{
			AppID:          appId,
			AppSecret:      appSecret,
			Token:          token,
			EncodingAESKey: encodingAESKey,
			Cache:          sysCache,
		})
		o.officialaccount = of
	} else {
		o.RUnlock()
	}
	return o.officialaccount, nil
}

func (o *oauth) GetWechat() (*WechatConfig, error) {
	reverseProxy, err := o.getOAuthConfig(OAuthConfigWechatReverseProxy)
	if err != nil {
		return nil, err
	}
	return &WechatConfig{
		Officialaccount: o.officialaccount,
		ReverseProxy:    reverseProxy,
	}, nil
}

func (o *oauth) WechatScanLoginRequest() (*dto.WechatScanLogin, error) {
	client, err := o.getWechatClient()
	if err != nil {
		return nil, err
	}
	code := utils.RandString(32, 10)
	ticket, err := client.GetBasic().GetQRTicket(&basic.Request{
		ExpireSeconds: 600,
		ActionName:    "QR_STR_SCENE",
		ActionInfo: struct {
			Scene struct {
				SceneStr string `json:"scene_str,omitempty"`
				SceneID  int    `json:"scene_id,omitempty"`
			} `json:"scene"`
		}{
			Scene: struct {
				SceneStr string `json:"scene_str,omitempty"`
				SceneID  int    `json:"scene_id,omitempty"`
			}{
				SceneStr: "login_" + code,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &dto.WechatScanLogin{
		URL:  "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + ticket.Ticket,
		Code: code,
	}, nil
}

func (o *oauth) WechatScanLogin(openid, code string) error {
	client, err := o.getWechatClient()
	if err != nil {
		return err
	}
	userinfo, err := client.GetUser().GetUserInfo(openid)
	if err != nil {
		return err
	}
	if userinfo.Subscribe == 0 {
		return errors.New("没有关注公众号")
	}
	wechat, err := o.wechatOAuthRepo.FindByOpenid(openid)
	if err != nil {
		return err
	}
	var uid int64
	if wechat == nil {
		// 新建账号
		user := &entity.User{
			Nickname:   userinfo.Nickname,
			Role:       "user",
			Createtime: time.Now().Unix(),
			Updatetime: time.Now().Unix(),
		}

		if uid, err = o.userSvc.oauthRegister(user); err != nil {
			return err
		} else {
			if userinfo.Headimgurl != "" {
				if b, err := util.HTTPGet(userinfo.Headimgurl); err == nil {
					if err := o.userSvc.UploadAvatar(uid, b); err != nil {
						logrus.Errorf("wechat register upload %s avatar: %v", userinfo.Headimgurl, err)
					}
				} else {
					logrus.Errorf("wechat register download %s avatar: %v", userinfo.Headimgurl, err)
				}
			}
		}
	} else {
		_, err := o.userSvc.UserInfo(wechat.UserID)
		if err != nil {
			return err
		}
		uid = wechat.UserID
	}
	if err := o.wechatOAuthRepo.BindCodeUid(code, uid); err != nil {
		return err
	}
	return nil
}

func (o *oauth) WechatScanLoginStatus(code string) (*dto.OAuthRespond, error) {
	uid, err := o.wechatOAuthRepo.FindCodeUid(code)
	if err != nil {
		return nil, err
	}
	user, err := o.userSvc.UserInfo(uid)
	if err != nil {
		return nil, err
	}
	return &dto.OAuthRespond{
		UserInfo: user,
		IsBind:   false,
	}, nil
}

func (o *oauth) getBbsClient() (*bbs.Client, error) {
	o.RLock()
	if o.wechat == nil {
		o.RUnlock()
		o.Lock()
		defer o.Unlock()
		clientid, err := o.getOAuthConfig(OAuthConfigBbsClientId)
		if err != nil {
			return nil, err
		}
		clientSecret, err := o.getOAuthConfig(OAuthConfigBbsClientSecret)
		if err != nil {
			return nil, err
		}
		o.bbs = bbs.NewClient(bbs.Config{
			ClientID:     clientid,
			ClientSecret: clientSecret,
		})
	} else {
		o.RUnlock()
	}
	return o.bbs, nil
}

func (o *oauth) getOAuthConfig(key string) (string, error) {
	ret, err := o.config.GetConfig(key)
	if err != nil {
		return "", err
	}
	if ret == "" {
		return "", errs.ErrOAuthPlatformNotConfigured
	}
	return ret, nil
}
