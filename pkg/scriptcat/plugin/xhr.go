package plugin

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dop251/goja"
	"github.com/scriptscat/cloudcat/pkg/scriptcat"
	"go.uber.org/zap"
)

func (g *GMPlugin) xmlHttpRequest(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) (func(call goja.FunctionCall) goja.Value, error) {
	cookieJar, err := g.gmFunc.LoadCookieJar(ctx, script)
	if err != nil {
		return nil, err
	}
	return func(call goja.FunctionCall) goja.Value {
		// TODO: 实现代理等
		cli := &http.Client{
			Transport: nil,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil
			},
			Jar:     cookieJar,
			Timeout: time.Second * 30,
		}

		if len(call.Arguments) != 1 {
			g.logger.Ctx(ctx).Warn("GMXHR incorrect number of parameters")
			return nil
		}
		args, ok := call.Arguments[0].Export().(map[string]interface{})
		if !ok {
			g.logger.Ctx(ctx).Warn("GMXHR parameter is not an object")
			return nil
		}
		method, _ := args["method"].(string)
		u, _ := args["url"].(string)
		if u == "" {
			g.logger.Ctx(ctx).Warn("GMXHR url cannot be empty")
			return nil
		}
		uu, err := url.Parse(u)
		if err != nil {
			g.logger.Ctx(ctx).Warn("GMXHR url format error", zap.Error(err))
			g.xhrOnError(ctx, runtime, args, err)
			return nil
		}
		uu.RawQuery = uu.Query().Encode()
		var body io.Reader
		if method != "GET" {
			data, _ := args["data"].(string)
			body = bytes.NewBufferString(data)
		}
		req, err := http.NewRequest(method, uu.String(), body)
		if err != nil {
			g.logger.Ctx(ctx).Warn("GMXHR create request failed", zap.Error(err))
			g.xhrOnError(ctx, runtime, args, err)
			return nil
		}
		// 默认header
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.40")
		req.Header.Set("Host", req.URL.Host)
		if headers, ok := args["headers"].(map[string]interface{}); ok {
			for k, v := range headers {
				req.Header.Set(k, v.(string))
			}
		}

		if timeout, _ := args["timeout"].(float64); timeout != 0 {
			cli.Timeout = time.Duration(timeout) * time.Millisecond
		}

		go func() {
			defer func(cookieJar CookieJar, ctx context.Context) {
				_ = cookieJar.Save(ctx)
			}(cookieJar, context.Background())
			resp, err := cli.Do(req)
			if err != nil {
				g.logger.Ctx(ctx).Warn("GMXHR request error", zap.Error(err))
				g.xhrOnError(ctx, runtime, args, err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				g.logger.Ctx(ctx).Warn("GMXHR request error", zap.Error(err))
				g.xhrOnError(ctx, runtime, args, err)
				return
			}
			respObj, err := goRespToXhrResp(resp, runtime, body)
			if err != nil {
				g.logger.Ctx(ctx).Warn("GMXHR request error", zap.Error(err))
				g.xhrOnError(ctx, runtime, args, err)
				return
			}
			onload, ok := goja.AssertFunction(runtime.ToValue(args["onload"]))
			if ok {
				g.logger.Ctx(ctx).Debug("GMXHR request onload")
				_, err := onload(nil, respObj)
				if err != nil {
					g.logger.Ctx(ctx).Warn("GMXHR onload error", zap.Error(err))
					return
				}
			}
		}()

		return goja.Undefined()
	}, nil
}

func goRespToXhrResp(resp *http.Response, runtime *goja.Runtime, body []byte) (goja.Value, error) {
	respObj := runtime.NewObject()
	if err := respObj.Set("finalUrl", resp.Request.URL.String()); err != nil {
		return nil, err
	}
	if err := respObj.Set("readyState", 4); err != nil {
		return nil, err
	}
	if err := respObj.Set("status", resp.StatusCode); err != nil {
		return nil, err
	}
	if err := respObj.Set("statusText", resp.Status); err != nil {
		return nil, err
	}
	respHeaders := ""
	for k, v := range resp.Header {
		respHeaders += k + ": " + v[0] + "\n"
	}
	if err := respObj.Set("responseHeaders", respHeaders); err != nil {
		return nil, err
	}
	if err := respObj.Set("response", string(body)); err != nil {
		return nil, err
	}
	if err := respObj.Set("responseText", string(body)); err != nil {
		return nil, err
	}
	return respObj, nil
}

func (g *GMPlugin) xhrOnError(ctx context.Context, runtime *goja.Runtime, args map[string]interface{}, err error) goja.Value {
	onerror, ok := goja.AssertFunction(runtime.ToValue(args["onerror"]))
	if ok {
		g.logger.Debug("GMXHR request onerror")
		errObj := runtime.NewObject()
		if err := errObj.Set("error", err.Error()); err != nil {
			return goja.Undefined()
		}
		_, err := onerror(nil, errObj)
		if err != nil {
			g.logger.Ctx(ctx).Warn("GMXHR onerror error", zap.Error(err))
			return goja.Undefined()
		}
	}
	return goja.Undefined()
}
