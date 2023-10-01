package cloudcat_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/pkg/utils"
)

type ConfigServer struct {
	BaseURL string `yaml:"baseURL"`
}

type ConfigUser struct {
	Name              string `yaml:"name"`
	Token             string `yaml:"token"`
	DataEncryptionKey string `yaml:"dataEncryptionKey"`
}

type Config struct {
	ApiVersion string        `yaml:"apiVersion"`
	Server     *ConfigServer `yaml:"server"`
	User       *ConfigUser   `yaml:"user"`
}

type Client struct {
	config *Config
	muxCli *mux.Client
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		muxCli: mux.NewClient(config.Server.BaseURL + "/api/" + config.ApiVersion),
	}
}

func (c *Client) Do(ctx context.Context, req any, resp any, opts ...mux.ClientOption) error {
	httpReq, err := c.muxCli.Request(ctx, req, opts...)
	if err != nil {
		return err
	}
	httpReq.ContentLength = 0
	httpReq.Header.Add("Authorization", "Bearer "+c.config.User.Token)
	encrypt, err := utils.NewAesEncrypt([]byte(c.config.User.DataEncryptionKey), httpReq.Body)
	if err != nil {
		return err
	}
	httpReq.Body = io.NopCloser(encrypt)
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return httpReq.Body, nil
	}
	// 请求
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()
	// warp body解密
	r, err := utils.NewAesDecrypt([]byte(c.config.User.DataEncryptionKey), httpResp.Body)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	jsonResp := &httputils.JSONResponse{
		Data: resp,
	}
	if err := json.Unmarshal(b, jsonResp); err != nil {
		return fmt.Errorf("json unmarshal error: %w", err)
	}
	if jsonResp.Code != 0 {
		return httputils.NewError(httpResp.StatusCode, jsonResp.Code, jsonResp.Msg)
	}
	return nil
}

var defaultClient *Client

func SetDefaultClient(cli *Client) {
	defaultClient = cli
}

func DefaultClient() *Client {
	return defaultClient
}
