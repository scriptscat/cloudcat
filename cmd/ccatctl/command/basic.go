package command

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/spf13/cobra"
)

type Basic struct {
	cli    *mux.Client
	config *string
	script *Script
	value  *Value
	cookie *Cookie
}

func NewBasic(config *string) *Basic {
	return &Basic{
		cli:    mux.NewClient("http://127.0.0.1:8080/api/v1"),
		config: config,
		script: NewScript(config),
		value:  NewValue(config),
		cookie: NewCookie(config),
	}
}

func (c *Basic) Command() []*cobra.Command {
	get := &cobra.Command{
		Use:   "get [resource]",
		Short: "获取资源信息",
	}
	get.AddCommand(c.script.Get(), c.value.Get(), c.cookie.Get())

	edit := &cobra.Command{
		Use:   "edit [resource]",
		Short: "编辑资源信息",
	}
	edit.AddCommand(c.script.Edit())

	cmd := []*cobra.Command{get, edit}
	cmd = append(cmd, c.script.Command()...)
	cmd = append(cmd, c.value.Command()...)
	cmd = append(cmd, c.cookie.Command()...)
	return cmd
}
