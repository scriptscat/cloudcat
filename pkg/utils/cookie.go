package utils

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/sirupsen/logrus"
)

type CookieFile struct {
	*cookiejar.Jar
}

func ReadCookie(cookie string) (http.CookieJar, error) {
	jar, _ := cookiejar.New(nil)
	m := make(map[string][]*http.Cookie)
	if err := json.Unmarshal([]byte(cookie), &m); err != nil && len(m) == 0 {
		return nil, err
	} else if err != nil {
		logrus.Errorf("cookie format error: %v", err)
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
			jar.SetCookies(u, []*http.Cookie{v})
		}
	}
	return jar, nil
}
