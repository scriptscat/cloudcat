package scripts_svc

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"

	"github.com/scriptscat/cloudcat/internal/model/entity/cookie_entity"
	"github.com/scriptscat/cloudcat/internal/model/entity/value_entity"
	"github.com/scriptscat/cloudcat/internal/repository/cookie_repo"
	"github.com/scriptscat/cloudcat/internal/repository/value_repo"
	"github.com/scriptscat/cloudcat/pkg/scriptcat"
	"github.com/scriptscat/cloudcat/pkg/scriptcat/plugin"
)

type GMPluginFunc struct {
}

func NewGMPluginFunc() plugin.GMPluginFunc {
	return &GMPluginFunc{}
}

func (g *GMPluginFunc) SetValue(ctx context.Context, script *scriptcat.Script, key string, value string) error {
	model, err := value_repo.Value().Find(ctx, script.StorageName(), key)
	if err != nil {
		return err
	}
	if model == nil {
		return value_repo.Value().Create(ctx, &value_entity.Value{
			StorageName: script.StorageName(),
			Key:         key,
			Value:       value,
			Createtime:  time.Now().Unix(),
		})
	}
	model.Value = value
	return value_repo.Value().Update(ctx, model)
}

func (g *GMPluginFunc) GetValue(ctx context.Context, script *scriptcat.Script, key string) (string, error) {
	model, err := value_repo.Value().Find(ctx, script.StorageName(), key)
	if err != nil {
		return "", err
	}
	if model == nil {
		return "", nil
	}
	return model.Value, nil
}

func (g *GMPluginFunc) ListValue(ctx context.Context, script *scriptcat.Script) (map[string]string, error) {
	list, _, err := value_repo.Value().FindPage(ctx, script.StorageName())
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, v := range list {
		m[v.Key] = v.Value
	}
	return m, nil
}

func (g *GMPluginFunc) DeleteValue(ctx context.Context, script *scriptcat.Script, key string) error {
	return value_repo.Value().Delete(ctx, script.StorageName(), key)
}

type cookieJar struct {
	sync.Mutex
	*cookiejar.Jar
	storageName string
	setUrl      map[string]*url.URL
}

func (c *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.Lock()
	defer c.Unlock()
	// 记录url
	c.setUrl[u.String()] = u
	// 设置cookie
	c.Jar.SetCookies(u, cookies)
}

func (c *cookieJar) Save(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()
	for u, v := range c.setUrl {
		cookies := c.Jar.Cookies(v)
		saveCookies := make([]*cookie_entity.HttpCookie, 0)
		for _, v := range cookies {
			saveCookies = append(saveCookies, &cookie_entity.HttpCookie{
				Name:       v.Name,
				Value:      v.Value,
				Path:       v.Path,
				Domain:     v.Domain,
				Expires:    v.Expires,
				RawExpires: v.RawExpires,
				MaxAge:     v.MaxAge,
				Secure:     v.Secure,
				HttpOnly:   v.HttpOnly,
				Raw:        v.Raw,
				Unparsed:   v.Unparsed,
			})
		}
		model, err := cookie_repo.Cookie().Find(ctx, c.storageName, u)
		if err != nil {
			return err
		}
		if model == nil {
			if err := cookie_repo.Cookie().Create(ctx, &cookie_entity.Cookie{
				StorageName: c.storageName,
				Url:         u,
				Cookies:     saveCookies,
				Createtime:  time.Now().Unix(),
			}); err != nil {
				return err
			}
		} else {
			if err := cookie_repo.Cookie().Update(ctx, &cookie_entity.Cookie{
				StorageName: c.storageName,
				Url:         u,
				Cookies:     saveCookies,
				Createtime:  model.Createtime,
			}); err != nil {
				return err
			}
		}
	}
	c.setUrl = make(map[string]*url.URL)
	return nil
}

func (g *GMPluginFunc) Logger(ctx context.Context, script *scriptcat.Script) *zap.Logger {
	return logger.Ctx(ctx).With(zap.String("script_id", script.ID),
		zap.String("name", script.Metadata["name"][0]))
}

func (g *GMPluginFunc) LoadCookieJar(ctx context.Context, script *scriptcat.Script) (plugin.CookieJar, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, err
	}
	cookies, _, err := cookie_repo.Cookie().FindPage(ctx, script.StorageName())
	if err != nil {
		return nil, err
	}
	for _, v := range cookies {
		u, err := url.Parse(v.Url)
		if err != nil {
			return nil, err
		}
		cookies := make([]*http.Cookie, 0)
		for _, v := range v.Cookies {
			cookies = append(cookies, &http.Cookie{
				Name:       v.Name,
				Value:      v.Value,
				Path:       v.Path,
				Domain:     v.Domain,
				Expires:    v.Expires,
				RawExpires: v.RawExpires,
				MaxAge:     v.MaxAge,
				Secure:     v.Secure,
				HttpOnly:   v.HttpOnly,
				Raw:        v.Raw,
				Unparsed:   v.Unparsed,
			})
		}
		jar.SetCookies(u, cookies)
	}
	return &cookieJar{
		Jar:         jar,
		storageName: script.StorageName(),
		setUrl:      make(map[string]*url.URL),
	}, nil
}
