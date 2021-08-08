package scriptcat

import (
	"rogchap.com/v8go"
)

func globalFunc(opts *Options, name string, callback v8go.FunctionCallback) {
	f, err := v8go.NewFunctionTemplate(opts.iso, callback)
	if err != nil {
		opts.err = err
		return
	}
	if err := opts.global.Set(name, f); err != nil {
		opts.err = err
	}
}

func GmXmlHttpRequest() Option {
	return func(opts *Options) {
		globalFunc(opts, "GM_xmlhttpRequest", func(info *v8go.FunctionCallbackInfo) *v8go.Value {

			return nil
		})
	}
}
