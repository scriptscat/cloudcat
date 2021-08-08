package scriptcat

import "rogchap.com/v8go"

type Context struct {
	ctx  *v8go.Context
	opts *Options
}

func NewContext(iso *Isolate, opt ...Option) (*Context, error) {
	global, err := v8go.NewObjectTemplate(iso.iso)
	if err != nil {
		return nil, err
	}
	options := &Options{
		iso:    iso.iso,
		global: global,
	}
	for _, o := range opt {
		o(options)
	}
	ctx, err := v8go.NewContext(iso.iso, options.global)
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
