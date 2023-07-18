package scripts_svc

import (
	"context"

	api "github.com/scriptscat/cloudcat/internal/api/scripts"
)

type CookieSvc interface {
	// CookieList 脚本cookie列表
	CookieList(ctx context.Context, req *api.CookieListRequest) (*api.CookieListResponse, error)
}

type cookieSvc struct {
}

var defaultCookie = &cookieSvc{}

func Cookie() CookieSvc {
	return defaultCookie
}

// CookieList 脚本cookie列表
func (c *cookieSvc) CookieList(ctx context.Context, req *api.CookieListRequest) (*api.CookieListResponse, error) {
	return nil, nil
}
