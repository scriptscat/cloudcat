package scriptcat

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/codfrm/cago/pkg/errs"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/dop251/goja"
	"go.uber.org/zap"
)

const (
	ScriptCat = "scriptcat"
)

// Runtime 脚本运行时
type Runtime interface {
	// Parse 解析脚本
	Parse(ctx context.Context, script string) (*Script, error)
	// Run 运行脚本
	Run(ctx context.Context, script *Script) (interface{}, error)
}

type scriptcat struct {
	plugins []Plugin
	logger  *logger.CtxLogger
}

var defaultRuntime Runtime

func RegisterRuntime(i Runtime) {
	defaultRuntime = i
}

func RuntimeCat() Runtime {
	return defaultRuntime
}

func NewRuntime(logger *logger.CtxLogger, plugins []Plugin) Runtime {
	return &scriptcat{
		plugins: plugins,
		logger:  logger,
	}
}

func (s *scriptcat) Name() string {
	return "scriptcat"
}

func (s *scriptcat) Parse(ctx context.Context, code string) (*Script, error) {
	meta := ParseMetaToJson(code)
	if len(meta["name"]) == 0 {
		return nil, errs.Warn(errors.New("script name is empty"))
	}
	if len(meta["namespace"]) == 0 {
		return nil, errs.Warn(errors.New("script namespace is empty"))
	}
	if len(meta["version"]) == 0 {
		return nil, errs.Warn(errors.New("script version is empty"))
	}
	id := meta["name"][0] + meta["namespace"][0]
	hash := sha256.New()
	hash.Write([]byte(id))
	id = fmt.Sprintf("%x", hash.Sum(nil))
	script := &Script{
		ID:       id,
		Code:     code,
		Metadata: meta,
	}
	return script, nil
}

func (s *scriptcat) Run(ctx context.Context, script *Script) (interface{}, error) {
	options := NewRunOptions()
	vm := goja.New()
	code := `
function vm` + script.ID + `() {
	` + script.Code + `
}
`
	for _, p := range s.plugins {
		if err := p.BeforeRun(ctx, script, vm); err != nil {
			return nil, err
		}
	}
	defer func() {
		for _, p := range s.plugins {
			if err := p.AfterRun(ctx, script, vm); err != nil {
				s.logger.Logger.Error("plugin after run error", zap.Error(err))
			}
		}
		vm.Interrupt("halt")
	}()
	_, err := vm.RunString(code)
	if err != nil {
		var e *goja.Exception
		if errors.As(err, &e) {
			s.logger.Ctx(ctx).Error("run script exception error",
				zap.String("error", e.Value().String()))
		} else {
			s.logger.Ctx(ctx).Error("run script error", zap.Error(err))
		}
		return nil, err
	}
	vmFun, ok := goja.AssertFunction(vm.Get("vm" + script.ID))
	if !ok {
		return nil, errors.New("not a vm function")
	}
	value, err := vmFun(goja.Undefined(), vm.ToValue(1), vm.ToValue(2))
	if err != nil {
		options.ResultCallback(nil, err)
		s.logger.Logger.Error("script run error", zap.Error(err))
		vm.Interrupt("halt")
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var result interface{}
	if err := s.then(vm, value, func() {
		result = value
		cancel()
	}); err != nil {
		s.logger.Logger.Error("script then error", zap.Error(err))
		return nil, err
	}
	if err := s.catch(vm, value, func() {
		cancel()
	}); err != nil {
		s.logger.Logger.Error("script catch error", zap.Error(err))
		return nil, err
	}
	<-ctx.Done()
	return result, nil
}

func (s *scriptcat) then(vm *goja.Runtime, value goja.Value, resolve func()) error {
	oj := value.ToObject(vm)
	then, ok := goja.AssertFunction(oj.Get("then"))
	if !ok {
		return errors.New("not a then function")
	}
	_, err := then(value, vm.ToValue(func(result interface{}) {
		// 任务完成处理
		s.logger.Logger.Info("script complete", zap.Any("result", result))
		resolve()
	}))
	if err != nil {
		return err
	}
	return nil
}

func (s *scriptcat) catch(vm *goja.Runtime, value goja.Value, reject func()) error {
	oj := value.ToObject(vm)
	catch, ok := goja.AssertFunction(oj.Get("catch"))
	if !ok {
		return errors.New("not a catch function")
	}
	_, err := catch(value, vm.ToValue(func(e interface{}) {
		promise, ok := oj.Export().(*goja.Promise)
		// 任务错误处理
		if ok {
			s.logger.Logger.Error("script error",
				zap.Any("error", e),
				zap.Any("reject", promise.Result().String()),
				//zap.Any("js stack",promise.Result())
			)
		} else {
			s.logger.Logger.Error("script error",
				zap.Any("error", e),
			)
		}
		reject()
	}))
	if err != nil {
		return err
	}
	return err
}
