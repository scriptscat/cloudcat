package scriptcat

import "net/http"

type Options struct {
	cookieJar http.CookieJar
}

type Option func(opts *Options)

func WithCookie(cookie http.CookieJar) Option {
	return func(opts *Options) {
		opts.cookieJar = cookie
	}
}
