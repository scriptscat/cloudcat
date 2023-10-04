package scripts

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/model/entity/cookie_entity"
)

type Cookie struct {
	StorageName string                      `json:"storage_name"`
	Host        string                      `json:"host"`
	Cookies     []*cookie_entity.HttpCookie `json:"cookies"`
	Createtime  int64                       `json:"createtime"`
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

// DeleteCookieRequest 删除cookie
type DeleteCookieRequest struct {
	mux.Meta    `path:"/cookies/:storageName" method:"DELETE"`
	StorageName string `uri:"storageName"`
	Host        string `form:"host"`
}

type DeleteCookieResponse struct {
}

// SetCookieRequest 设置cookie
type SetCookieRequest struct {
	mux.Meta    `path:"/cookies/:storageName" method:"POST"`
	StorageName string                      `uri:"storageName"`
	Cookies     []*cookie_entity.HttpCookie `form:"cookies"`
}

// SetCookieResponse 设置cookie
type SetCookieResponse struct {
}
