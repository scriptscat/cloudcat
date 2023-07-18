package scripts

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/pkg/scriptcat/cookie"
)

type Cookie struct {
	CookieSpace string           `json:"cookie_space"`
	Url         string           `json:"url"`
	Cookies     []*cookie.Cookie `json:"cookies"`
	CreatedAt   int64            `json:"created_at"`
}

// CookieListRequest 脚本cookie列表
type CookieListRequest struct {
	mux.Meta `path:"/cookies" method:"GET"`
}

// CookieListResponse 脚本cookie列表
type CookieListResponse struct {
	List []*Cookie `json:"list"`
}
