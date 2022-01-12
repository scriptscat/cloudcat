package executor

import (
	"github.com/sirupsen/logrus"
	"rogchap.com/v8go"
)

func Console() Option {
	return func(opts *Options) {
		if opts._log == nil {
			return
		}
		console, _ := v8go.NewObjectTemplate(opts.iso)
		logfn, _ := v8go.NewFunctionTemplate(opts.iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			for _, v := range info.Args() {
				if v.IsObject() {
					if b, err := v.MarshalJSON(); err != nil {
						opts.log(logrus.InfoLevel, "%s", v.String())
					} else {
						opts.log(logrus.InfoLevel, "%s", b)
					}
				} else {
					opts.log(logrus.InfoLevel, v.DetailString())
				}
			}
			return nil
		})
		if err := console.Set("log", logfn); err != nil {
			opts.err = err
			return
		}

		if conObj, err := console.NewInstance(opts.ctx); err != nil {
			opts.err = err
		} else if err := opts.ctx.Global().Set("console", conObj); err != nil {
			opts.err = err
		}
	}
}
