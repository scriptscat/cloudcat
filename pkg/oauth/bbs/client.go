package bbs

import (
	"encoding/json"
	"net/url"

	"github.com/scriptscat/cloudcat/pkg/utils"
)

const serverUrl = "https://bbs.tampermonkey.net.cn"

type Config struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
}

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) RequestAccessToken(code string) (*AccessTokenRespond, error) {
	resp, err := utils.HttpPost(serverUrl+"/plugin.php?id=codfrm_oauth2:server&op=access_token", "client_id="+c.config.ClientID+"&client_secret="+
		c.config.ClientSecret+"&code="+code, nil)
	if err != nil {
		return nil, err
	}
	ret := &AccessTokenRespond{}
	if err := json.Unmarshal(resp, ret); err != nil {
		return nil, err
	}
	if ret.Code != 0 {
		return nil, ret
	}
	return ret, nil
}

func (c *Client) RequestUser(access_token string) (*UserRespond, error) {
	resp, err := utils.HttpPost(serverUrl+"/plugin.php?id=codfrm_oauth2:server&op=user", "access_token="+access_token, nil)
	if err != nil {
		return nil, err
	}
	ret := &UserRespond{}
	if err := json.Unmarshal(resp, ret); err != nil {
		return nil, err
	}
	if ret.Code != 0 {
		return nil, ret
	}
	return ret, nil
}

func (c *Client) RedirectURL(redirectUrl string) string {
	return serverUrl + "/plugin.php?id=codfrm_oauth2:oauth&client_id=" + c.config.ClientID +
		"&scope=user&response_type=code&redirect_uri=" + url.PathEscape(redirectUrl)
}
