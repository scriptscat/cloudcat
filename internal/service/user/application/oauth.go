package application

import (
	"errors"
	"sync"
	"time"

	config2 "github.com/scriptscat/cloudcat/internal/infrastructure/config"
	"github.com/scriptscat/cloudcat/internal/pkg/kvdb"
	entity2 "github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/errs"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/repository"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/vo"
	"github.com/scriptscat/cloudcat/pkg/cnt"
	"github.com/scriptscat/cloudcat/pkg/oauth/bbs"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	offOAuth "github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/silenceper/wechat/v2/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OAuth interface {
	RedirectOAuth(redirectUrl, platform string) (string, error)
	BBSOAuthLogin(code string) (*vo.OAuthRespond, error)
	WechatAuthLogin(code string) (*vo.OAuthRespond, error)

	WechatScanRequest(op string) (*vo.WechatScan, error)
	WechatScanLogin(openid, code string) error
	WechatScanLoginStatus(code string) (*vo.OAuthRespond, error)

	GetWechat() (*WechatConfig, error)

	OAuthPlatform(uid int64) (*vo.OpenPlatform, error)

	WechatScanBind(openid, code string) error
	WechatScanBindCode(uid int64, code string) error

	BindBbs(uid int64, code string) error
	Unbind(uid int64, platform string) error
}

type WechatConfig struct {
	Officialaccount *officialaccount.OfficialAccount
	ReverseProxy    string
}

type oauth struct {
	sync.RWMutex
	userSvc         User
	kv              kvdb.KvDb
	config          config2.SystemConfig
	bbsOAuthRepo    repository.BBSOAuth
	wechatOAuthRepo repository.WechatOAuth
	tx              *gorm.DB

	wechat          *offOAuth.Oauth
	bbs             *bbs.Client
	officialaccount *officialaccount.OfficialAccount
}

func NewOAuth(config config2.SystemConfig, kv kvdb.KvDb, tx *gorm.DB, userSvc User, bbs repository.BBSOAuth, wc repository.WechatOAuth) OAuth {
	return &oauth{
		config:          config,
		bbsOAuthRepo:    bbs,
		wechatOAuthRepo: wc,
		kv:              kv,
		tx:              tx,
		userSvc:         userSvc,
	}
}

func (o *oauth) RedirectOAuth(redirectUrl, platform string) (string, error) {
	homeUrl, _ := o.config.GetConfig(config2.HomeUrl)
	redirectUrl = homeUrl + redirectUrl
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

func (o *oauth) BBSOAuthLogin(code string) (*vo.OAuthRespond, error) {
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
	var uid int64
	if bbs == nil {
		// 需要绑定账号登录
		user := &entity2.User{
			Username:   userResp.User.Username,
			Role:       "user",
			Createtime: time.Now().Unix(),
			Updatetime: time.Now().Unix(),
		}
		if err := o.userSvc.CheckUsername(user.Username); err != nil {
			if err != errs.ErrUsernameExist {
				return nil, err
			}
			user.Username += utils.RandString(4, 2)
		}
		if err := o.tx.Transaction(func(tx *gorm.DB) error {
			uid, err = o.userSvc.oauthRegister(tx, user)
			if err != nil {
				return err
			}
			return repository.NewBbsOAuth(tx).Save(&entity2.BbsOauthUser{
				Openid:     userResp.User.Uid,
				UserID:     uid,
				Status:     cnt.ACTIVE,
				Createtime: time.Now().Unix(),
			})
		}); err != nil {
			return nil, err
		}
		if userResp.User.Avatar != "" {
			if b, err := util.HTTPGet(userResp.User.Avatar); err == nil {
				if err := o.userSvc.UploadAvatar(uid, b); err != nil {
					logrus.Errorf("bbs register upload %s avatar: %v", userResp.User.Avatar, err)
				}
			} else {
				logrus.Errorf("bbs register download %s avatar: %v", userResp.User.Avatar, err)
			}
		}
	} else {
		uid = bbs.UserID
	}
	user, err := o.userSvc.UserInfo(uid)
	if err != nil {
		return nil, err
	}
	return &vo.OAuthRespond{
		UserInfo: user,
		IsBind:   true,
	}, nil
}

func (o *oauth) WechatAuthLogin(code string) (*vo.OAuthRespond, error) {
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
		return &vo.OAuthRespond{
			IsBind: false,
		}, nil
	}
	user, err := o.userSvc.UserInfo(wechat.UserID)
	if err != nil {
		return nil, err
	}
	return &vo.OAuthRespond{
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
		appId, err := o.getOAuthConfig(config2.OAuthConfigWechatAppId)
		if err != nil {
			return nil, err
		}
		appSecret, err := o.getOAuthConfig(config2.OAuthConfigWechatAppSecret)
		if err != nil {
			return nil, err
		}
		token, err := o.getOAuthConfig(config2.OAuthConfigWechatToken)
		if err != nil {
			return nil, err
		}
		encodingAESKey, err := o.getOAuthConfig(config2.OAuthConfigWechatEncodingaeskey)
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
	reverseProxy, err := o.getOAuthConfig(config2.OAuthConfigWechatReverseProxy)
	if err != nil {
		return nil, err
	}
	wc, err := o.getWechatClient()
	if err != nil {
		return nil, err
	}
	return &WechatConfig{
		Officialaccount: wc,
		ReverseProxy:    reverseProxy,
	}, nil
}

func (o *oauth) WechatScanRequest(op string) (*vo.WechatScan, error) {
	client, err := o.getWechatClient()
	if err != nil {
		return nil, err
	}
	code := utils.RandString(16, 1)
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
				SceneStr: op + "_" + code,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &vo.WechatScan{
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
		user := &entity2.User{
			Username:   userinfo.Nickname,
			Role:       "user",
			Createtime: time.Now().Unix(),
			Updatetime: time.Now().Unix(),
		}
		if err := o.userSvc.CheckUsername(user.Username); err != nil {
			if err != errs.ErrUsernameExist {
				return err
			}
			user.Username += utils.RandString(4, 2)
		}
		if err := o.tx.Transaction(func(tx *gorm.DB) error {
			uid, err = o.userSvc.oauthRegister(tx, user)
			if err != nil {
				return err
			}
			return repository.NewWechatOAuth(tx, o.kv).Save(&entity2.WechatOauthUser{
				Openid:     openid,
				Unionid:    userinfo.UnionID,
				UserID:     uid,
				Status:     cnt.ACTIVE,
				Createtime: time.Now().Unix(),
			})
		}); err != nil {
			return err
		}
		if userinfo.Headimgurl != "" {
			if b, err := util.HTTPGet(userinfo.Headimgurl); err == nil {
				if err := o.userSvc.UploadAvatar(uid, b); err != nil {
					logrus.Errorf("wechat register upload %s avatar: %v", userinfo.Headimgurl, err)
				}
			} else {
				logrus.Errorf("wechat register download %s avatar: %v", userinfo.Headimgurl, err)
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
	return client.GetCustomerMessageManager().Send(&message.CustomerMessage{
		ToUser:  openid,
		Msgtype: "text",
		Text: &message.MediaText{
			Content: "扫码登录成功",
		},
	})
}

func (o *oauth) WechatScanLoginStatus(code string) (*vo.OAuthRespond, error) {
	uid, err := o.wechatOAuthRepo.FindCodeUid(code)
	if err != nil {
		return nil, err
	}
	if uid == 0 {
		return nil, errs.ErrRecordNotFound
	}
	user, err := o.userSvc.UserInfo(uid)
	if err != nil {
		return nil, err
	}
	return &vo.OAuthRespond{
		UserInfo: user,
		IsBind:   true,
	}, nil
}

func (o *oauth) getBbsClient() (*bbs.Client, error) {
	o.RLock()
	if o.bbs == nil {
		o.RUnlock()
		o.Lock()
		defer o.Unlock()
		clientid, err := o.getOAuthConfig(config2.OAuthConfigBbsClientId)
		if err != nil {
			return nil, err
		}
		clientSecret, err := o.getOAuthConfig(config2.OAuthConfigBbsClientSecret)
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

func (o *oauth) OAuthPlatform(uid int64) (*vo.OpenPlatform, error) {
	bbs, err := o.bbsOAuthRepo.FindByUid(uid)
	if err != nil {
		return nil, err
	}
	wechat, err := o.wechatOAuthRepo.FindByUid(uid)
	if err != nil {
		return nil, err
	}
	return &vo.OpenPlatform{
		Bbs:    bbs != nil,
		Wechat: wechat != nil,
	}, nil
}

// NOTE:code的利用方法和登录相反

func (o *oauth) WechatScanBind(openid, code string) error {
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
	// 消费掉这个code
	uid, err := o.wechatOAuthRepo.FindCodeUid(code)
	if err != nil {
		return err
	}
	if uid == 0 {
		return errs.ErrRecordNotFound
	}
	if wechat != nil {
		return client.GetCustomerMessageManager().Send(&message.CustomerMessage{
			ToUser:  openid,
			Msgtype: "text",
			Text: &message.MediaText{
				Content: "该微信已绑定过账号了",
			},
		})
	}
	w, err := o.wechatOAuthRepo.FindByUid(uid)
	if err != nil {
		return err
	}
	if w != nil {
		return client.GetCustomerMessageManager().Send(&message.CustomerMessage{
			ToUser:  openid,
			Msgtype: "text",
			Text: &message.MediaText{
				Content: "该账号已经绑定过其他微信了",
			},
		})
	}
	if err := o.wechatOAuthRepo.Save(&entity2.WechatOauthUser{
		Openid:     openid,
		UserID:     uid,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
	}); err != nil {
		return err
	}
	return client.GetCustomerMessageManager().Send(&message.CustomerMessage{
		ToUser:  openid,
		Msgtype: "text",
		Text: &message.MediaText{
			Content: "绑定成功",
		},
	})
}

func (o *oauth) WechatScanBindCode(uid int64, code string) error {
	return o.wechatOAuthRepo.BindCodeUid(code, uid)
}

func (o *oauth) BindBbs(uid int64, code string) error {
	client, err := o.getBbsClient()
	if err != nil {
		return err
	}
	resp, err := client.RequestAccessToken(code)
	if err != nil {
		return err
	}
	userResp, err := client.RequestUser(resp.AccessToken)
	if err != nil {
		return err
	}
	bbs, err := o.bbsOAuthRepo.FindByOpenid(userResp.User.Uid)
	if err != nil {
		return err
	}
	if bbs != nil {
		return errs.ErrBindOtherUser
	}
	if u, err := o.bbsOAuthRepo.FindByUid(uid); err != nil {
		return err
	} else if u != nil {
		return errs.ErrBindOtherOAuth
	}
	return o.bbsOAuthRepo.Save(&entity2.BbsOauthUser{
		Openid:     userResp.User.Uid,
		UserID:     uid,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
	})
}

func (o *oauth) Unbind(uid int64, platform string) error {
	switch platform {
	case "bbs":
		bbs, err := o.bbsOAuthRepo.FindByUid(uid)
		if err != nil {
			return err
		}
		if bbs.Createtime+2592000 > time.Now().Unix() {
			return errs.ErrNotUnbind
		}
		return o.bbsOAuthRepo.Delete(bbs.ID)
	case "wechat":
		wechat, err := o.wechatOAuthRepo.FindByUid(uid)
		if err != nil {
			return err
		}
		if wechat.Createtime+2592000 > time.Now().Unix() {
			return errs.ErrNotUnbind
		}
		return o.wechatOAuthRepo.Delete(wechat.ID)
	}
	return nil
}
