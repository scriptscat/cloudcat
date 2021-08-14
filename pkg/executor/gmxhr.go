package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"rogchap.com/v8go"
)

func globalFunc(opts *Options, name string, callback v8go.FunctionCallback) {
	f, err := v8go.NewFunctionTemplate(opts.iso, callback)
	if err != nil {
		opts.err = err
		return
	}
	if err := opts.ctx.Global().Set(name, f.GetFunction(opts.ctx)); err != nil {
		opts.err = err
	}
}

func GmXmlHttpRequest(jar http.CookieJar) Option {
	return func(opts *Options) {
		globalFunc(opts, "GM_xmlhttpRequest", func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			//TODO: 实现代理等
			cli := &http.Client{
				Transport:     nil,
				CheckRedirect: nil,
				Jar:           jar,
				Timeout:       time.Second * 30,
			}
			if len(info.Args()) != 1 {
				opts.log(logrus.ErrorLevel, "GMXHR number of parameters")
				return nil
			}
			arg := info.Args()[0]
			if !arg.IsObject() {
				opts.log(logrus.ErrorLevel, "GMXHR arg not object")
				return nil
			}
			details := arg.Object()
			var err error
			defer func() {
				if err != nil {
					if onerror := getFunction(details, "onerror"); onerror != nil {
						_, _ = onerror.Call()
					}
				}
			}()
			method := strings.ToUpper(getObjString(details, "method"))
			u := getObjString(details, "url")
			var data io.Reader
			if method != "GET" {
				data = bytes.NewBufferString(getObjString(details, "data"))
			}
			req, err := http.NewRequest(method, u, data)
			if err != nil {
				opts.log(logrus.ErrorLevel, "GMXHR New Request: %v", err)
				return nil
			}
			if cookie := getObjString(details, "cookie"); cookie != "" {
				if anonymous := getObjBool(details, "anonymous"); !anonymous && jar != nil {
					u, _ := url.Parse(u)
					cookies := jar.Cookies(u)
					for _, v := range cookies {
						cookie = fmt.Sprintf("%v=%v;", v.Name, v.Value) + cookie
					}
				}
				req.Header.Add("Cookie", cookie)
			}

			if headers := getObject(details, "headers"); headers != nil {
				headersMap := make(map[string]string)
				b, _ := headers.MarshalJSON()
				if err := json.Unmarshal(b, &headersMap); err == nil {
					for k, v := range headersMap {
						req.Header.Set(k, v)
					}
				}
			}

			if timeout := getNumber(details, "timeout"); timeout != 0 {
				cli.Timeout = time.Duration(timeout) * time.Millisecond
			}

			go func() {
				var err error
				defer func() {
					if err != nil {
						if onerror := getFunction(details, "onerror"); onerror != nil {
							_, _ = onerror.Call()
						}
					}
				}()
				resp, err := cli.Do(req)
				if err != nil {
					opts.log(logrus.ErrorLevel, "GMXHR Request: %v", err)
					if err == http.ErrHandlerTimeout {
						if ontimeout := getFunction(details, "ontimeout"); ontimeout != nil {
							_, _ = ontimeout.Call()
						}
						return
					}
					return
				}

				// 处理resp
				if onload := getFunction(details, "onload"); onload != nil {
					var xhrResp *v8go.ObjectTemplate
					xhrResp, err = v8go.NewObjectTemplate(opts.iso)
					if err != nil {
						opts.log(logrus.ErrorLevel, "GMXHR Respond: %v", err)
						if onerror := getFunction(details, "onerror"); onerror != nil {
							_, _ = onerror.Call()
						}
						return
					}
					xhrResp.Set("finalUrl", u)
					xhrResp.Set("readyState", 4)
					responseHeaders := ""
					for k, v := range resp.Header {
						for _, v := range v {
							responseHeaders = responseHeaders + fmt.Sprintf("%s: %s\n", k, v)
						}
					}
					xhrResp.Set("responseHeaders", responseHeaders)
					var body []byte
					body, err = io.ReadAll(resp.Body)
					xhrResp.Set("status", resp.StatusCode)
					xhrResp.Set("responseText", string(body))

					var arg *v8go.Object
					if arg, err = xhrResp.NewInstance(info.Context()); err != nil {
						opts.log(logrus.ErrorLevel, "GMXHR Respond: %v", err)
						return
					}

					if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
						if v, err := v8go.JSONParse(info.Context(), string(body)); err != nil {
							arg.Set("response", string(body))
						} else {
							arg.Set("response", v)
						}
					} else {
						arg.Set("response", string(body))
					}

					onload.Call(arg)
				}

			}()
			return nil
		})
	}
}
