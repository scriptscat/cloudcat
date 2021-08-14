package scriptcat

import "github.com/scriptscat/cloudcat/pkg/executor"

type ScriptCat struct {
	script string
	meta   map[string][]string
	code   string
	opts   *Options
}

func NewScriptCat() (*ScriptCat, error) {

	return &ScriptCat{}, nil
}

func (s *ScriptCat) Run(script string, opt ...Option) error {
	ctx, code, err := s.compile(script, opt...)
	if err != nil {
		return err
	}
	_, err = ctx.RunScript(code, "main.js")
	return err
}

func (s *ScriptCat) RunOnce(script string, opt ...Option) error {

	return nil
}

func (s *ScriptCat) compile(script string, opt ...Option) (*executor.Context, string, error) {
	options := &Options{}
	for _, o := range opt {
		o(options)
	}
	exec, err := executor.NewExecutor()
	if err != nil {
		return nil, "", err
	}
	// 解析script
	metaJson := ParseMetaToJson(ParseMeta(script))
	ctx, err := s.buildContext(exec, metaJson, options)
	if err != nil {
		return nil, "", err
	}
	// TODO: 编译code(require resource等内容)
	return ctx, "", nil
}

func (s *ScriptCat) buildContext(exec *executor.Executor, meta map[string][]string, opts *Options) (*executor.Context, error) {
	contextOpts := make([]executor.Option, 0)

	optMap := map[string]func() executor.Option{
		"GM_xmlhttpRequest": func() executor.Option {
			return executor.GmXmlHttpRequest(opts.cookieJar)
		},
	}

	for _, v := range meta["grant"] {
		if f, ok := optMap[v]; ok {
			contextOpts = append(contextOpts, f())
		}
	}

	ctx, err := executor.NewContext(exec, contextOpts...)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}
