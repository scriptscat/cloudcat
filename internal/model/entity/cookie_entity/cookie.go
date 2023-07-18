package cookie_entity

import (
	"time"

	"github.com/scriptscat/cloudcat/pkg/scriptcat/cookie"
)

type Cookie struct {
	CookieSpace string           `json:"cookie_space"`
	Url         string           `json:"url"`
	Cookies     []*cookie.Cookie `json:"cookies"`
	CreatedTime time.Time        `json:"created_time"`
}
