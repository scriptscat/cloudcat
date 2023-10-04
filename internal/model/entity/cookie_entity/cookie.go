package cookie_entity

import (
	"fmt"
	"net/http"
	"time"
)

type HttpCookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Path    string    `json:"path"`    // optional
	Domain  string    `json:"domain"`  // optional
	Expires time.Time `json:"expires"` // optional

	ExpirationDate int64 `json:"expirationDate,omitempty"` // optional 到期时间戳

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int    `json:"max_age"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"http_only"`
	SameSite string `json:"same_site"`
}

type Cookie struct {
	StorageName string        `json:"storage_name"`
	Host        string        `json:"host"`
	Cookies     []*HttpCookie `json:"cookies"`
	Createtime  int64         `json:"createtime"`
}

func (h *HttpCookie) ToCookie(cookie *http.Cookie) {
	h.Name = cookie.Name
	h.Value = cookie.Value
	h.Path = cookie.Path
	h.Domain = cookie.Domain
	h.Expires = cookie.Expires
	h.MaxAge = cookie.MaxAge
	h.Secure = cookie.Secure
	h.HttpOnly = cookie.HttpOnly
	//switch cookie.SameSite {
	//case http.SameSiteDefaultMode:
	//	h.SameSite = "default"
	//case http.SameSiteLaxMode:
	//	h.SameSite = "lax"
	//case http.SameSiteStrictMode:
	//	h.SameSite = "strict"
	//case http.SameSiteNoneMode:
	//	h.SameSite = "none"
	//default:
	//	h.SameSite = ""
	//}
}

func (h *HttpCookie) ToHttpCookie() *http.Cookie {
	return &http.Cookie{
		Name:     h.Name,
		Value:    h.Value,
		Path:     h.Path,
		Domain:   h.Domain,
		Expires:  h.Expires,
		MaxAge:   h.MaxAge,
		Secure:   h.Secure,
		HttpOnly: h.HttpOnly,
		//SameSite: http.SameSiteDefaultMode,
	}
}

// ID returns the domain;path;name triple of e as an ID.
func (e *HttpCookie) ID() string {
	return fmt.Sprintf("%s;%s;%s", e.Domain, e.Path, e.Name)
}
