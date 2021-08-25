package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type CookieFile struct {
	*cookiejar.Jar
}

type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

func ReadCookie(cookie string) (http.CookieJar, error) {
	jar, _ := cookiejar.New(nil)
	m := make(map[string][]*Cookie)
	if err := json.Unmarshal([]byte(cookie), &m); err != nil && len(m) == 0 {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("cookie format error: %v", err)
	}
	for k, v := range m {
		for _, v := range v {
			urlStr := ""
			if v.Secure {
				urlStr = "https://"
			} else {
				urlStr = "http://"
			}
			urlStr += v.Domain
			u, _ := url.Parse("https://" + k)
			jar.SetCookies(u, []*http.Cookie{{
				Name:       v.Name,
				Value:      v.Value,
				Path:       v.Path,
				Domain:     v.Domain,
				Expires:    v.Expires,
				RawExpires: v.RawExpires,
				MaxAge:     v.MaxAge,
				Secure:     v.Secure,
				HttpOnly:   v.HttpOnly,
				Raw:        v.Raw,
				Unparsed:   v.Unparsed,
			}})
		}
	}
	return jar, nil
}
