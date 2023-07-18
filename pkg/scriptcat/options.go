package scriptcat

type RunOptions struct {
	resultCallback func(result interface{}, err error)
}

type RunOption func(*RunOptions)

func NewRunOptions(opts ...RunOption) *RunOptions {
	options := &RunOptions{}
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithResultCallback(callback func(result interface{}, err error)) RunOption {
	return func(options *RunOptions) {
		options.resultCallback = callback
	}
}

func (o *RunOptions) ResultCallback(result interface{}, err error) {
	if o.resultCallback != nil {
		o.resultCallback(result, err)
	}
}
