package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"github.com/scriptscat/cloudcat/pkg/scriptcat"
)

func (g *GMPlugin) setValue(ctx context.Context, script *scriptcat.Script, runtime *goja.Runtime) (func(call goja.FunctionCall) goja.Value, error) {
	return func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		arg1 := call.Argument(1)
		export := arg1.Export()
		value, err := json.Marshal(export)
		if err != nil {
			panic(fmt.Errorf("GM_setValue error: %v", err))
		}

		if err := g.gmFunc.SetValue(ctx, script, key, string(value)); err != nil {
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
		if s == "" {
			if len(call.Arguments) > 1 {
				return call.Argument(1)
			}
			return goja.Undefined()
		}
		var v interface{}
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			panic(fmt.Errorf("GM_getValue error: %v", err))
		}
		return runtime.ToValue(v)
	}, nil
}
