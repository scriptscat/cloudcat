package plugin

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
	"github.com/scriptscat/cloudcat/pkg/scriptcat"
)

func (g *GMPlugin) setValue(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) (func(call goja.FunctionCall) goja.Value, error) {
	return func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		arg1 := call.Argument(1)
		if err := g.gmFunc.SetValue(ctx, script, key, arg1.Export()); err != nil {
			panic(fmt.Errorf("GM_setValue error: %v", err))
		}
		return goja.Undefined()
	}, nil
}

func (g *GMPlugin) getValue(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) (func(call goja.FunctionCall) goja.Value, error) {
	return func(call goja.FunctionCall) goja.Value {
		s, err := g.gmFunc.GetValue(ctx, script, call.Argument(0).String())
		if err != nil {
			return goja.Undefined()
		}
		if s == nil {
			if len(call.Arguments) > 1 {
				return call.Argument(1)
			}
			return goja.Undefined()
		}
		return runtime.ToValue(s)
	}, nil
}
