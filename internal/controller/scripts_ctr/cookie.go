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
