package scripts_svc

import (
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/scriptscat/cloudcat/internal/model/entity/cookie_entity"
	"github.com/scriptscat/cloudcat/internal/model/entity/resource_entity"
	"github.com/scriptscat/cloudcat/internal/model/entity/value_entity"
	"github.com/scriptscat/cloudcat/internal/repository/cookie_repo"
	"github.com/scriptscat/cloudcat/internal/repository/resource_repo"
	"github.com/scriptscat/cloudcat/internal/repository/value_repo"
	"github.com/scriptscat/cloudcat/pkg/scriptcat"
	"github.com/scriptscat/cloudcat/pkg/scriptcat/plugin"
	"go.uber.org/zap"
)

type GMPluginFunc struct {
}

func NewGMPluginFunc() plugin.GMPluginFunc {
	return &GMPluginFunc{}
}

func (g *GMPluginFunc) SetValue(ctx context.Context, script *scriptcat.Script, key string, value interface{}) error {
	model, err := value_repo.Value().Find(ctx, script.StorageName(), key)
	if err != nil {
		return err
	}
	if model == nil {
		model = &value_entity.Value{
			StorageName: script.StorageName(),
			Key:         key,
			Createtime:  time.Now().Unix(),
		}
		err := model.Value.Set(value)
		if err != nil {
			return err
		}
		return value_repo.Value().Create(ctx, model)
	}
	if err := model.Value.Set(value); err != nil {
		return err
	}
	return value_repo.Value().Update(ctx, model)
}

func (g *GMPluginFunc) GetValue(ctx context.Context, script *scriptcat.Script, key string) (interface{}, error) {
	model, err := value_repo.Value().Find(ctx, script.StorageName(), key)
	if err != nil {
		return "", err
	}
	if model == nil {
		return "", nil
	}
	return model.Value.Get(), nil
}

func (g *GMPluginFunc) ListValue(ctx context.Context, script *scriptcat.Script) (map[string]interface{}, error) {
	list, _, err := value_repo.Value().FindPage(ctx, script.StorageName())
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	for _, v := range list {
		m[v.Key] = v.Value.Get()
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
	cookies     map[string]map[string]*cookie_entity.HttpCookie
}

func (c *cookieJar) Cookies(u *url.URL) []*http.Cookie {
	cookies := c.Jar.Cookies(u)
	return cookies
}

func (c *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.Lock()
	defer c.Unlock()
	// 记录url
	host, err := cookie_entity.CanonicalHost(u.Host)
	if err != nil {
		return
	}
	key := cookie_entity.JarKey(host, nil)
	defPath := cookie_entity.DefaultPath(u.Path)
	submap, ok := c.cookies[key]
	if !ok {
		submap = make(map[string]*cookie_entity.HttpCookie)
	}
	for _, v := range cookies {
		cookie := &cookie_entity.HttpCookie{}
		cookie.ToCookie(v)
		if cookie.Path == "" || cookie.Path[0] != '/' {
			cookie.Path = defPath
		}
		if cookie.Domain != "" && cookie.Domain[0] != '.' {
			cookie.Domain = host
		}
		submap[cookie.ID()] = cookie
	}
	if len(submap) == 0 {
		delete(c.cookies, key)
	} else {
		c.cookies[key] = submap
	}
	// 设置cookie
	c.Jar.SetCookies(u, cookies)
}

func (c *cookieJar) Save(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()
	for host, v := range c.cookies {
		saveCookies := make([]*cookie_entity.HttpCookie, 0)
		for _, v := range v {
			saveCookies = append(saveCookies, v)
		}
		model, err := cookie_repo.Cookie().Find(ctx, c.storageName, host)
		if err != nil {
			return err
		}
		if model == nil {
			if err := cookie_repo.Cookie().Create(ctx, &cookie_entity.Cookie{
				StorageName: c.storageName,
				Host:        host,
				Cookies:     saveCookies,
				Createtime:  time.Now().Unix(),
			}); err != nil {
				return err
			}
		} else {
			model.Cookies = mergeCookie(model.Cookies, saveCookies)
			if err := cookie_repo.Cookie().Update(ctx, model); err != nil {
				return err
			}
		}
	}
	c.cookies = make(map[string]map[string]*cookie_entity.HttpCookie)
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
		u, err := url.Parse("https://" + v.Host)
		if err != nil {
			return nil, err
		}
		cookies := make([]*http.Cookie, 0)
		for _, v := range v.Cookies {
			if v.Domain != "" && v.Domain[0] != '.' {
				u, err := url.Parse("https://" + v.Domain)
				if err != nil {
					return nil, err
				}
				jar.SetCookies(u, []*http.Cookie{v.ToHttpCookie()})
				continue
			}
			cookies = append(cookies, v.ToHttpCookie())
		}
		jar.SetCookies(u, cookies)
	}
	return &cookieJar{
		Jar:         jar,
		storageName: script.StorageName(),
		cookies:     make(map[string]map[string]*cookie_entity.HttpCookie),
	}, nil
}

func (g *GMPluginFunc) LoadResource(ctx context.Context, url string) (string, error) {
	// 从资源缓存中搜索是否存在
	resource, err := resource_repo.Resource().Find(ctx, url)
	if err != nil {
		return "", err
	}
	if resource != nil {
		return resource.Content, nil
	}
	// 从远程获取资源
	resp, err := http.Get(url) // #nosec
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 保存资源
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	m := &resource_entity.Resource{
		URL:        url,
		Content:    string(content),
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	if err := resource_repo.Resource().Create(ctx, m); err != nil {
		return "", err
	}
	return m.Content, nil
}
