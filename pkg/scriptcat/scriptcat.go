package scriptcat

import (
	"errors"
	"sync"

	"github.com/scriptscat/cloudcat/pkg/executor"
	"github.com/sirupsen/logrus"
	"rogchap.com/v8go"
)

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
	ctx, code, err := s.compile(script, opt...)
	if err != nil {
		return err
	}
	ret, err := ctx.RunScript(code, "main.js")
	if err != nil {
		return err
	}
	if !ret.IsPromise() {
		return errors.New("return is not a promise object")
	}
	l := sync.WaitGroup{}
	p, err := ret.AsPromise()
	if err != nil {
		return err
	}
	l.Add(1)
	p.Then(func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		l.Done()
		return nil
	})
	p.Catch(func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		l.Done()
		return nil
	})
	l.Wait()
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

	return ctx, "function main() {\n" + script + "\n}\nmain();", nil
}

func (s *ScriptCat) buildContext(exec *executor.Executor, meta map[string][]string, opts *Options) (*executor.Context, error) {
	contextOpts := []executor.Option{
		executor.WithLogger(logrus.StandardLogger().Logf),
		executor.Console(),
	}

	optMap := map[string]func() executor.Option{
		"GM_xmlhttpRequest": func() executor.Option {
			return executor.GmXmlHttpRequest(opts.cookieJar)
		},
		"GM_notification": func() executor.Option {
			return executor.GmNotification()
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
