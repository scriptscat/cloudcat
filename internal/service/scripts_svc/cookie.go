package scripts_svc

import (
	"context"

	"github.com/codfrm/cago/pkg/i18n"
	api "github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/internal/pkg/code"
	"github.com/scriptscat/cloudcat/internal/repository/cookie_repo"
	"github.com/scriptscat/cloudcat/internal/repository/script_repo"
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
	scripts, err := script_repo.Script().FindByStorage(ctx, req.StorageName)
	if err != nil {
		return nil, err
	}
	if len(scripts) == 0 {
		return nil, i18n.NewNotFoundError(ctx, code.StorageNameNotFound)
	}
	script := scripts[0]
	list, _, err := cookie_repo.Cookie().FindPage(ctx, script.StorageName())
	if err != nil {
		return nil, err
	}
	resp := &api.CookieListResponse{
		List: make([]*api.Cookie, 0),
	}
	for _, v := range list {
		resp.List = append(resp.List, &api.Cookie{
			StorageName: v.StorageName,
			Url:         v.Url,
			Cookies:     v.Cookies,
			Createtime:  v.Createtime,
		})
	}
	return resp, nil
}
