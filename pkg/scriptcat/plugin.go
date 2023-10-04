package scriptcat

import (
	"context"

	"github.com/dop251/goja"
)

type Plugin interface {
	// BeforeRun 运行前
	BeforeRun(ctx context.Context, script *Script, runtime *goja.Runtime) error
	// AfterRun 运行后
	AfterRun(ctx context.Context, script *Script, runtime *goja.Runtime) error
}
