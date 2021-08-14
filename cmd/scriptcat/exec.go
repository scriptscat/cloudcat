package main

import (
	"io/ioutil"

	"github.com/scriptscat/cloudcat/pkg/scriptcat"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
)

type execCmd struct {
	cookiefile string
}

func newExecCmd() *execCmd {
	return &execCmd{}
}

func (e *execCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "exec [file] [flags]",
		Short: "执行一个脚本猫脚本",
		RunE:  e.exec,
	}
	ret.Flags().StringVarP(&e.cookiefile, "cookiefile", "c", "", "设置cookie文件")

	return []*cobra.Command{ret}
}

func (e *execCmd) exec(cmd *cobra.Command, args []string) error {

	sc, err := scriptcat.NewScriptCat()
	if err != nil {
		return err
	}

	script, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	opts := make([]scriptcat.Option, 0)
	if e.cookiefile != "" {
		jar, err := utils.ReadCookie(e.cookiefile)
		if err != nil {
			return err
		}
		opts = append(opts, scriptcat.WithCookie(jar))
	}

	return sc.Run(string(script))
}
