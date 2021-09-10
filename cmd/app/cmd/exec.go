package cmd

import (
	"archive/zip"
	"context"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path"

	"github.com/scriptscat/cloudcat/pkg/scriptcat"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type execCmd struct {
	cookiefile string
	runOnce    bool
}

func NewExecCmd() *execCmd {
	return &execCmd{}
}

func (e *execCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "exec [file] [flags]",
		Short: "执行一个脚本猫脚本",
		RunE:  e.exec,
		Args:  cobra.ExactArgs(1),
	}
	ret.Flags().StringVarP(&e.cookiefile, "cookiefile", "c", "", "设置cookie文件")
	ret.Flags().BoolVarP(&e.runOnce, "run-once", "", false, "让定时脚本只运行一次")

	return []*cobra.Command{ret}
}

func (e *execCmd) exec(cmd *cobra.Command, args []string) error {
	var err error
	var script, cookie, value fs.File
	if path.Ext(args[0]) != ".js" {
		// 软件包
		pkg, err := zip.OpenReader(args[0])
		if err != nil {
			return err
		}
		script, err = pkg.Open("userScript.js")
		if err != nil {
			return err
		}
		defer script.Close()
		cookie, _ = pkg.Open("cookie.json")
		value, _ = pkg.Open("value.json")
	} else {
		script, err = os.Open(args[0])
		if err != nil {
			return err
		}
		defer script.Close()
		cookie, err = os.Open(e.cookiefile)
		if err != nil {
			return err
		}
		defer cookie.Close()
	}

	opts := []scriptcat.Option{scriptcat.WithLogger(logrus.StandardLogger().Logf)}
	if cookie != nil {
		jar, err := utils.ReadCookie(readString(cookie))
		if err != nil {
			return err
		}
		opts = append(opts, scriptcat.WithCookie(jar))
	}

	if value != nil {
		opts = append(opts, scriptcat.WithValue(value))
	}

	sc, err := scriptcat.NewScriptCat()
	if err != nil {
		return err
	}

	if e.runOnce {
		_, err = sc.RunOnce(context.Background(), readString(script), opts...)
	} else {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		_, err = sc.Run(ctx, readString(script), opts...)
	}
	return err
}

func readString(r io.Reader) string {
	if r == nil {
		return ""
	}
	byte, _ := io.ReadAll(r)
	return string(byte)
}
