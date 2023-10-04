package scripts_ctr

import (
	"context"

	api "github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/internal/service/scripts_svc"
)

type Cookie struct {
}

func NewCookie() *Cookie {
	return &Cookie{}
}

// CookieList 脚本cookie列表
func (c *Cookie) CookieList(ctx context.Context, req *api.CookieListRequest) (*api.CookieListResponse, error) {
	return scripts_svc.Cookie().CookieList(ctx, req)
}

// DeleteCookie 删除cookie
func (c *Cookie) DeleteCookie(ctx context.Context, req *api.DeleteCookieRequest) (*api.DeleteCookieResponse, error) {
	return scripts_svc.Cookie().DeleteCookie(ctx, req)
}

// SetCookie 设置cookie
func (c *Cookie) SetCookie(ctx context.Context, req *api.SetCookieRequest) (*api.SetCookieResponse, error) {
	return scripts_svc.Cookie().SetCookie(ctx, req)
}
