package cookie_entity

import (
	"github.com/scriptscat/cloudcat/pkg/scriptcat/cookie"
)

type Cookie struct {
	StorageName string           `json:"storage_name"`
	Url         string           `json:"url"`
	Cookies     []*cookie.Cookie `json:"cookies"`
	Createtime  int64            `json:"createtime"`
}
