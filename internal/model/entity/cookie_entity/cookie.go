package cookie_entity

import (
	"time"
)

type HttpCookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Path       string    `json:"path"`        // optional
	Domain     string    `json:"domain"`      // optional
	Expires    time.Time `json:"expires"`     // optional
	RawExpires string    `json:"raw_expires"` // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int      `json:"max_age"`
	Secure   bool     `json:"secure"`
	HttpOnly bool     `json:"http_only"`
	SameSite int      `json:"same_site"`
	Raw      string   `json:"raw"`
	Unparsed []string `json:"unparsed"` // Raw text of unparsed attribute-value pairs
}

type Cookie struct {
	StorageName string        `json:"storage_name"`
	Url         string        `json:"url"`
	Cookies     []*HttpCookie `json:"cookies"`
	Createtime  int64         `json:"createtime"`
}
