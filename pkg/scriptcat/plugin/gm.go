package plugin

import (
	"context"
	"errors"
	"net/http"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/dop251/goja"
	"github.com/goccy/go-json"
	"github.com/scriptscat/cloudcat/pkg/scriptcat"
	"go.uber.org/zap"
)

type CookieJar interface {
	http.CookieJar
	Save(ctx context.Context) error
}

type GMPluginFunc interface {
	SetValue(ctx context.Context, script *scriptcat.Script, key string, value interface{}) error
	GetValue(ctx context.Context, script *scriptcat.Script, key string) (interface{}, error)
	ListValue(ctx context.Context, script *scriptcat.Script) (map[string]interface{}, error)
	DeleteValue(ctx context.Context, script *scriptcat.Script, key string) error

	Logger(ctx context.Context, script *scriptcat.Script) *zap.Logger

	LoadCookieJar(ctx context.Context, script *scriptcat.Script) (CookieJar, error)

	LoadResource(ctx context.Context, url string) (string, error)
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
		logger: logger.NewCtxLogger(logger.Default()).With(zap.String("plugin", "GMPlugin")),
		gmFunc: storage,
		ctxMap: make(map[string]*ctxCancel),
	}
	p.grantMap = map[string]grantFunc{
		"GM_xmlhttpRequest": p.xmlHttpRequest,
		"GM_setValue":       p.setValue,
		"GM_getValue":       p.getValue,
		"GM_log":            p.log,
		"GM_notification":   p.empty,
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
	// 注入require
	for _, v := range script.Metadata["require"] {
		s, err := g.gmFunc.LoadResource(ctx, v)
		if err != nil {
			logger.Ctx(ctx).Error("load resource error", zap.Error(err))
		} else {
			_, err := runtime.RunString(s)
			if err != nil {
				var e *goja.Exception
				if errors.As(err, &e) {
					logger.Ctx(ctx).Error("run script exception error",
						zap.String("error", e.Value().String()))
				} else {
					logger.Ctx(ctx).Error("run script error", zap.Error(err))
				}
				return err
			}
		}
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

func (g *GMPlugin) empty(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) (func(call goja.FunctionCall) goja.Value, error) {
	return func(call goja.FunctionCall) goja.Value {
		logger.Ctx(ctx).Debug("empty function")
		return goja.Undefined()
	}, nil
}
