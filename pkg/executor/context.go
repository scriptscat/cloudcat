package executor

import "rogchap.com/v8go"

type Context struct {
	ctx  *v8go.Context
	opts *Options
}

func NewContext(executor *Executor, opt ...Option) (*Context, error) {
	global, err := v8go.NewObjectTemplate(executor.iso)
	if err != nil {
		return nil, err
	}
	options := &Options{
		iso:    executor.iso,
		global: global,
	}
	for _, o := range opt {
		o(options)
	}
	ctx, err := v8go.NewContext(executor.iso, options.global)
	if err != nil {
		return nil, err
	}
	return &Context{
		ctx: ctx,
	}, nil
}

func (c *Context) RunScript(source string, origin string) (*v8go.Value, error) {
	return c.ctx.RunScript(source, origin)
}
