package scriptcat

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/scriptscat/cloudcat/pkg/executor"
	"github.com/sirupsen/logrus"
	"rogchap.com/v8go"
)

type ScriptCat struct {
	script string
	meta   map[string][]string
	code   string
}

func NewScriptCat() (*ScriptCat, error) {

	return &ScriptCat{}, nil
}

func (s *ScriptCat) options(opt ...Option) *Options {
	options := &Options{
		location: time.Local,
	}
	for _, o := range opt {
		o(options)
	}
	return options
}

func (s *ScriptCat) Run(ctx context.Context, script string, opt ...Option) (string, error) {
	opts := s.options(opt...)
	exec, meta, code, err := s.compile(script, opts)
	if err != nil {
		return "", err
	}
	// 判断是否是定时脚本
	c, ok := meta["crontab"]
	if !ok {
		opts.log(logrus.InfoLevel, "start run script")
		return s.runOnce(ctx, exec, code)
	}
	cron := cron.New(cron.WithSeconds(), cron.WithLocation(opts.location))
	unit := strings.Split(c[0], " ")
	if len(unit) == 5 {
		unit = append([]string{"0"}, unit...)
	}
	// 对once进行处理
	for i, v := range unit {
		if v == "once" {
			unit[i] = "*"
			i -= 1
			for ; i >= 0; i-- {
				if unit[i] == "*" {
					unit[i] = "0"
				}
			}
			break
		}
	}
	c[0] = strings.Join(unit, " ")
	opts.log(logrus.InfoLevel, "start run crontab script: %s", c[0])
	_, err = cron.AddFunc(c[0], func() {
		ret, err := s.runOnce(ctx, exec, code)
		if err != nil {
			opts.log(logrus.ErrorLevel, "run script error: %v", err)
		} else {
			opts.log(logrus.InfoLevel, "run script ok: %v", ret)
		}
	})
	if err != nil {
		return "", err
	}
	cron.Start()
	<-ctx.Done()
	return "", nil
}

func (s *ScriptCat) RunOnce(ctx context.Context, script string, opt ...Option) (string, error) {
	opts := s.options(opt...)
	exec, _, code, err := s.compile(script, opts)
	if err != nil {
		return "", err
	}
	ret, err := s.runOnce(ctx, exec, code)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (s *ScriptCat) runOnce(ctx context.Context, exec *executor.Context, code string) (msg string, err error) {
	ret, err := exec.RunScript(code, "app.js")
	if err != nil {
		return "", err
	}
	if !ret.IsPromise() {
		return "", errors.New("return is not a promise object")
	}
	p, err := ret.AsPromise()
	if err != nil {
		return "", err
	}
	done := make(chan struct{})
	p.Then(func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		if len(info.Args()) == 1 && info.Args()[0].IsString() {
			msg = info.Args()[0].String()
		}
		done <- struct{}{}
		return nil
	})
	p.Catch(func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		if len(info.Args()) == 1 && info.Args()[0].IsString() {
			err = errors.New(info.Args()[0].String())
		}
		done <- struct{}{}
		return nil
	})
	select {
	case <-done:
	case <-ctx.Done():
	}
	return
}

func (s *ScriptCat) compile(script string, options *Options) (*executor.Context, map[string][]string, string, error) {
	exec, err := executor.NewExecutor()
	if err != nil {
		return nil, nil, "", err
	}
	// 解析script
	metaJson := ParseMetaToJson(ParseMeta(script))
	ctx, err := s.buildContext(exec, metaJson, options)
	if err != nil {
		return nil, nil, "", err
	}
	// TODO: 编译code(require resource等内容)

	return ctx, metaJson, "function app() {\n" + script + "\n}\napp();", nil
}

func (s *ScriptCat) buildContext(exec *executor.Executor, meta map[string][]string, opts *Options) (*executor.Context, error) {
	contextOpts := []executor.Option{
		executor.WithLogger(opts.log),
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
