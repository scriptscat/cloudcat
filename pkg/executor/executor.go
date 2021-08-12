package executor

import "rogchap.com/v8go"

type Executor struct {
	iso *v8go.Isolate
}

func NewExecutor() (*Executor, error) {
	iso, err := v8go.NewIsolate()
	if err != nil {
		return nil, err
	}
	return &Executor{iso: iso}, nil
}

func (c *Executor) Run(ctx *Context, source string) (*v8go.Value, error) {
	return ctx.RunScript(source, "main.js")
}
