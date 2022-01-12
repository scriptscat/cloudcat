package executor

import "rogchap.com/v8go"

type Executor struct {
	iso *v8go.Isolate
	ctx *v8go.Context
}

func NewExecutor(opt ...Option) (*Executor, error) {
	iso, err := v8go.NewIsolate()
	if err != nil {
		return nil, err
	}
	global, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, err
	}
	options := &Options{
		iso:    iso,
		global: global,
	}
	ctx, err := v8go.NewContext(iso, options.global)
	if err != nil {
		return nil, err
	}
	options.ctx = ctx
	for _, o := range opt {
		o(options)
		if options.err != nil {
			return nil, err
		}
	}
	return &Executor{
		iso: iso,
		ctx: ctx,
	}, nil
}

func (c *Executor) Run(source string) (*v8go.Value, error) {
	return c.ctx.RunScript(source, "app.js")
}
