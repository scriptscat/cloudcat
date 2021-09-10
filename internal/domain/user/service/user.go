package service

import (
	"github.com/scriptscat/cloudcat/internal/domain/user/dto"
	"github.com/scriptscat/cloudcat/internal/domain/user/errs"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/pkg/oauth/bbs"
)

type User interface {
	RedirectOAuth(redirectUrl, platform string) (string, error)
	BBSOAuthLogin(code string) (*dto.UserInfo, error)
}

const (
	OAUTH_CONFIG_BBS_CLIENT_ID     = "oauth_config_bbs_client_id"
	OAUTH_CONFIG_BBS_CLIENT_SECRET = "oauth_config_bbs_client_secret"
)

type user struct {
	config config.SystemConfig
}

func NewUser(config config.SystemConfig) User {
	return &user{config: config}
}

func (u *user) RedirectOAuth(redirectUrl, platform string) (string, error) {
	switch platform {
	case "bbs":
		client, err := u.getBbsClient()
		if err != nil {
			return "", err
		}
		return client.RedirectURL(redirectUrl), err
	}
	return "", errs.ErrOAuthPlatformNotSupport
}

func (u *user) BBSOAuthLogin(code string) (*dto.UserInfo, error) {
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

	return nil, nil
}

func (u *user) getBbsClient() (*bbs.Client, error) {
	clientid, err := u.config.GetConfig(OAUTH_CONFIG_BBS_CLIENT_ID)
	if err != nil {
		return nil, err
	}
	if clientid == "" {
		return nil, errs.ErrOAuthPlatformNotConfigured
	}
	clientSecret, err := u.config.GetConfig(OAUTH_CONFIG_BBS_CLIENT_SECRET)
	if err != nil {
		return nil, err
	}
	if clientid == "" {
		return nil, errs.ErrOAuthPlatformNotConfigured
	}
	return bbs.NewClient(bbs.Config{
		ClientID:     clientid,
		ClientSecret: clientSecret,
	}), nil
}
