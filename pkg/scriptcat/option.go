package scriptcat

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Options struct {
	context.Context
	cookieJar http.CookieJar
	_log      func(level logrus.Level, format string, args ...interface{})
}

type Option func(opts *Options)

func WithCookie(cookie http.CookieJar) Option {
	return func(opts *Options) {
		opts.cookieJar = cookie
	}
}

func WithValue(value interface{}) Option {
	return func(opts *Options) {

	}
}

func WithLogger(log func(level logrus.Level, format string, args ...interface{})) Option {
	return func(opts *Options) {
		opts._log = log
	}
}

func (o *Options) log(level logrus.Level, format string, args ...interface{}) {
	if o._log != nil {
		o._log(level, format, args)
	}
}
