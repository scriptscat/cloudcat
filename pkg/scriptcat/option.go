package scriptcat

import "rogchap.com/v8go"

type Options struct {
	iso    *v8go.Isolate
	global *v8go.ObjectTemplate
	err    error
}

type Option func(opts *Options)
