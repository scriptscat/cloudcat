package scripts

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/pkg/scriptcat/cookie"
)

type Cookie struct {
	StorageName string           `json:"storage_name"`
	Url         string           `json:"url"`
	Cookies     []*cookie.Cookie `json:"cookies"`
	Createtime  int64            `json:"createtime"`
}

// CookieListRequest 脚本cookie列表
type CookieListRequest struct {
	mux.Meta    `path:"/cookies/:storageName" method:"GET"`
	StorageName string `uri:"storageName"`
}

// CookieListResponse 脚本cookie列表
type CookieListResponse struct {
	List []*Cookie `json:"list"`
}
