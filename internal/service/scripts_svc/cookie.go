package scripts_svc

import (
	"context"
	"time"

	"github.com/codfrm/cago/pkg/i18n"
	api "github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/internal/model/entity/cookie_entity"
	"github.com/scriptscat/cloudcat/internal/pkg/code"
	"github.com/scriptscat/cloudcat/internal/repository/cookie_repo"
	"github.com/scriptscat/cloudcat/internal/repository/script_repo"
)

type CookieSvc interface {
	// CookieList 脚本cookie列表
	CookieList(ctx context.Context, req *api.CookieListRequest) (*api.CookieListResponse, error)
	// DeleteCookie 删除cookie
	DeleteCookie(ctx context.Context, req *api.DeleteCookieRequest) (*api.DeleteCookieResponse, error)
	// SetCookie 设置cookie
	SetCookie(ctx context.Context, req *api.SetCookieRequest) (*api.SetCookieResponse, error)
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
			Host:        v.Host,
			Cookies:     v.Cookies,
			Createtime:  v.Createtime,
		})
	}
	return resp, nil
}

// DeleteCookie 删除cookie
func (c *cookieSvc) DeleteCookie(ctx context.Context, req *api.DeleteCookieRequest) (*api.DeleteCookieResponse, error) {
	scripts, err := script_repo.Script().FindByStorage(ctx, req.StorageName)
	if err != nil {
		return nil, err
	}
	if len(scripts) == 0 {
		return nil, i18n.NewNotFoundError(ctx, code.StorageNameNotFound)
	}
	if err := cookie_repo.Cookie().Delete(ctx, scripts[0].StorageName(), req.Host); err != nil {
		return nil, err
	}
	return nil, nil
}

// SetCookie 设置cookie
func (c *cookieSvc) SetCookie(ctx context.Context, req *api.SetCookieRequest) (*api.SetCookieResponse, error) {
	scripts, err := script_repo.Script().FindByStorage(ctx, req.StorageName)
	if err != nil {
		return nil, err
	}
	if len(scripts) == 0 {
		return nil, i18n.NewNotFoundError(ctx, code.StorageNameNotFound)
	}
	cookiesMap := make(map[string][]*cookie_entity.HttpCookie)
	for _, v := range req.Cookies {
		host, err := cookie_entity.CanonicalHost(v.Domain)
		if err != nil {
			return nil, err
		}
		key := cookie_entity.JarKey(host, nil)
		_, ok := cookiesMap[key]
		if !ok {
			cookiesMap = make(map[string][]*cookie_entity.HttpCookie)
		}
		cookiesMap[key] = append(cookiesMap[key], v)
	}
	for key, v := range cookiesMap {
		model, err := cookie_repo.Cookie().Find(ctx, scripts[0].StorageName(), key)
		if err != nil {
			return nil, err
		}
		if model == nil {
			if err := cookie_repo.Cookie().Create(ctx, &cookie_entity.Cookie{
				StorageName: scripts[0].StorageName(),
				Host:        key,
				Cookies:     v,
				Createtime:  time.Now().Unix(),
			}); err != nil {
				return nil, err
			}
		} else {
			model.Cookies = mergeCookie(model.Cookies, v)
			if err := cookie_repo.Cookie().Update(ctx, model); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func mergeCookie(oldCookie, newCookie []*cookie_entity.HttpCookie) []*cookie_entity.HttpCookie {
	cookieMap := make(map[string]*cookie_entity.HttpCookie)
	for _, v := range oldCookie {
		cookieMap[v.ID()] = v
	}
	for _, v := range newCookie {
		cookieMap[v.ID()] = v
	}
	cookies := make([]*cookie_entity.HttpCookie, 0)
	now := time.Now()
	for _, v := range cookieMap {
		if v.MaxAge < 0 {
			continue
		} else if v.MaxAge > 0 {
			if v.Expires.IsZero() {
				v.Expires = now.Add(time.Duration(v.MaxAge) * time.Second)
			}
		} else {
			if !v.Expires.IsZero() {
				if !v.Expires.After(now) {
					continue
				}
			}
		}
		cookies = append(cookies, v)
	}
	return cookies
}
