package executor

import (
	"github.com/sirupsen/logrus"
	"rogchap.com/v8go"
)

type Options struct {
	iso    *v8go.Isolate
	global *v8go.ObjectTemplate
	ctx    *v8go.Context
	err    error
	_log   func(level logrus.Level, format string, args ...interface{})
}

type Option func(opts *Options)

func WithLogger(log func(level logrus.Level, format string, args ...interface{})) Option {
	return func(opts *Options) {
		opts._log = log
	}
}

func (o *Options) log(level logrus.Level, format string, args ...interface{}) {
	if o._log != nil {
		o._log(level, format, args...)
	}
}
