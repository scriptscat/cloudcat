package main

import (
	"github.com/scriptscat/cloudcat/pkg/scriptcat"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type execCmd struct {
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

	return sc.Run(string(script))
}
