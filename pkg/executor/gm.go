package executor

import (
	"github.com/sirupsen/logrus"
	"rogchap.com/v8go"
)

func GmNotification() Option {
	return func(opts *Options) {
		globalFunc(opts, "GM_notification", func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if len(info.Args()) == 0 {
				return nil
			}
			var title, text string
			if info.Args()[0].IsObject() {
				arg1, err := info.Args()[0].AsObject()
				if err != nil {
					return nil
				}
				title = getObjString(arg1, "title")
				text = getObjString(arg1, "text")
				if fn := getFunction(arg1, "ondone"); fn != nil {
					fn.Call()
				}
				if len(info.Args()) == 2 && info.Args()[1].IsFunction() {
					fn, err := info.Args()[0].AsFunction()
					if err != nil {
						return nil
					}
					fn.Call()
				}
			} else {
				switch len(info.Args()) {
				case 2:
					text = getString(info.Args()[1])
					fallthrough
				case 1:
					title = getString(info.Args()[0])
				}
			}
			opts.log(logrus.InfoLevel, "Notification: %v %v", title, text)

			return nil
		})
	}
}
