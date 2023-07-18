package window

import (
	"context"

	"github.com/dop251/goja"
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
	return nil
}

func (w *Plugin) AfterRun(ctx context.Context, script *scriptcat2.Script, runtime *goja.Runtime) error {
	return nil
}
