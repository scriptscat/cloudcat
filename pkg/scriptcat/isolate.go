package scriptcat

import "rogchap.com/v8go"

type Isolate struct {
	iso *v8go.Isolate
}

func NewIsolate() (*Isolate, error) {
	iso, err := v8go.NewIsolate()
	if err != nil {
		return nil, err
	}
	return &Isolate{iso: iso}, nil
}

func (c *Isolate) Run(ctx *Context, source string) (*v8go.Value, error) {
	return ctx.RunScript(source, "main.js")
}
