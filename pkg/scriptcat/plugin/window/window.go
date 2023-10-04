package window

import (
	"context"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/dop251/goja_nodejs/url"
	scriptcat2 "github.com/scriptscat/cloudcat/pkg/scriptcat"
)

// Plugin 注册一些基础的函数
type Plugin struct {
}

func NewBrowserPlugin() scriptcat2.Plugin {
	return &Plugin{}
}

func (w *Plugin) Name() string {
	return "NodeJS"
}

func (w *Plugin) BeforeRun(ctx context.Context, script *scriptcat2.Script, vm *goja.Runtime) error {
	// 注册计时器
	timer := NewTimer(vm)
	if err := timer.Start(); err != nil {
		return err
	}
	if err := vm.Set("window", vm.GlobalObject()); err != nil {
		return err
	}
	// url
	nodejs := new(require.Registry)
	nodejs.Enable(vm)
	url.Enable(vm)
	return nil
}

func (w *Plugin) AfterRun(ctx context.Context, script *scriptcat2.Script, runtime *goja.Runtime) error {
	return nil
}
