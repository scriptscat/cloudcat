package plugin

import (
	"bytes"
	"context"
	"github.com/codfrm/cago/pkg/logger"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"go.uber.org/zap"

	"github.com/dop251/goja"
	"github.com/scriptscat/cloudcat/pkg/scriptcat"
)

type CookieJar interface {
	http.CookieJar
	Save(ctx context.Context) error
}

type GMPluginFunc interface {
	SetValue(ctx context.Context, script *scriptcat.Script, key string, value string) error
	GetValue(ctx context.Context, script *scriptcat.Script, key string) (string, error)
	ListValue(ctx context.Context, script *scriptcat.Script) (map[string]string, error)
	DeleteValue(ctx context.Context, script *scriptcat.Script, key string) error

	Logger(ctx context.Context, script *scriptcat.Script) *zap.Logger

	LoadCookieJar(ctx context.Context, script *scriptcat.Script) (CookieJar, error)
}

type grantFunc func(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) (func(call goja.FunctionCall) goja.Value, error)

type ctxCancel struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// GMPlugin gm 函数插件
type GMPlugin struct {
	logger   *logger.CtxLogger
	gmFunc   GMPluginFunc
	grantMap map[string]grantFunc
	ctxMap   map[string]*ctxCancel
}

func NewGMPlugin(storage GMPluginFunc) scriptcat.Plugin {
	p := &GMPlugin{
		logger: logger.NewCtxLogger(logger.Default()),
		gmFunc: storage,
		ctxMap: make(map[string]*ctxCancel),
	}
	p.grantMap = map[string]grantFunc{
		"GM_xmlhttpRequest": p.xmlHttpRequest,
		"GM_setValue":       p.setValue,
		"GM_getValue":       p.getValue,
		"GM_log":            p.log,
	}
	return p
}

func (g *GMPlugin) Name() string {
	return "GMPlugin"
}

func (g *GMPlugin) Version() string {
	return "0.0.1"
}

func (g *GMPlugin) BeforeRun(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) error {
	// 根据meta注入gm函数
	ctx, cancel := context.WithCancel(ctx)
	g.ctxMap[script.ID] = &ctxCancel{
		ctx:    ctx,
		cancel: cancel,
	}
	// 默认注入GM_log
	defaultGrant := []string{"GM_log"}
	for _, v := range defaultGrant {
		f, err := g.grantMap[v](ctx, script, runtime)
		if err != nil {
			return err
		}
		if err := runtime.Set(v, f); err != nil {
			return err
		}
	}
	for _, grant := range script.Metadata["grant"] {
		g, ok := g.grantMap[grant]
		if !ok {
			continue
		}
		f, err := g(ctx, script, runtime)
		if err != nil {
			return err
		}
		if err := runtime.Set(grant, f); err != nil {
			return err
		}
	}
	return nil
}

func (g *GMPlugin) AfterRun(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) error {
	if v, ok := g.ctxMap[script.ID]; ok {
		v.cancel()
		delete(g.ctxMap, script.ID)
	}
	return nil
}

func (g *GMPlugin) log(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) (func(call goja.FunctionCall) goja.Value, error) {
	return func(call goja.FunctionCall) goja.Value {
		msg := ""
		level := "info"
		labels := make([]zap.Field, 0)
		if len(call.Arguments) >= 1 {
			msg = call.Argument(0).String()
			if len(call.Arguments) >= 2 {
				level = call.Argument(1).String()
				if len(call.Arguments) >= 3 {
					b, err := json.Marshal(call.Argument(2))
					if err == nil {
						labels = append(labels, zap.ByteString("labels", b))
					}
				}
			}
		}

		logger := g.gmFunc.Logger(ctx, script)

		switch level {
		case "debug":
			logger.Debug(msg, labels...)
		case "info":
			logger.Info(msg, labels...)
		case "warn":
			logger.Warn(msg, labels...)
		case "error":
			logger.Error(msg, labels...)
		}

		return goja.Undefined()
	}, nil
}

func (g *GMPlugin) xmlHttpRequest(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) (func(call goja.FunctionCall) goja.Value, error) {
	cookieJar, err := g.gmFunc.LoadCookieJar(ctx, script)
	if err != nil {
		return nil, err
	}
	return func(call goja.FunctionCall) goja.Value {
		// TODO: 实现代理等
		cli := &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           cookieJar,
			Timeout:       time.Second * 30,
		}

		if len(call.Arguments) != 1 {
			g.logger.Warn("GMXHR 参数数量不正确")
			return nil
		}
		arg, ok := call.Arguments[0].Export().(map[string]interface{})
		if !ok {
			g.logger.Warn("GMXHR 参数不是对象")
			return nil
		}
		method, _ := arg["method"].(string)
		url, _ := arg["url"].(string)
		if url == "" {
			g.logger.Warn("GMXHR url不能为空")
			return nil
		}
		var body io.Reader
		if method != "GET" {
			data, _ := arg["data"].(string)
			body = bytes.NewBufferString(data)
		}
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			g.logger.Warn("GMXHR 创建请求失败", zap.Error(err))
			return nil
		}
		if headers, ok := arg["headers"].(map[string]interface{}); ok {
			for k, v := range headers {
				req.Header.Set(k, v.(string))
			}
		}

		if timeout, _ := arg["timeout"].(float64); timeout != 0 {
			cli.Timeout = time.Duration(timeout) * time.Millisecond
		}

		go func() {
			defer cookieJar.Save(context.Background())
			resp, err := cli.Do(req)
			if err != nil {
				g.logger.Warn("GMXHR 请求失败", zap.Error(err))
				return
			}
			defer resp.Body.Close()
		}()

		return goja.Undefined()
	}, nil
}
