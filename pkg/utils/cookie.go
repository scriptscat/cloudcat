package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type CookieFile struct {
	*cookiejar.Jar
}

func ReadCookie(filename string) (http.CookieJar, error) {
	jar, _ := cookiejar.New(nil)
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	m := make(map[string][]*http.Cookie)
	if err := json.Unmarshal(f, &m); err != nil {
		return nil, err
	}
	for k, v := range m {
		u, err := url.Parse(k)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", k, err)
		}
		jar.SetCookies(u, v)
	}
	return jar, nil
}
