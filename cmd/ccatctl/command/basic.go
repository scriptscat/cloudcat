package command

import (
	"github.com/spf13/cobra"
)

type Basic struct {
	script *Script
	value  *Value
	cookie *Cookie
	token  *Token
}

func NewBasic(config string) *Basic {
	return &Basic{
		script: NewScript(),
		value:  NewValue(),
		cookie: NewCookie(),
		token:  NewToken(),
	}
}

func (c *Basic) Command() []*cobra.Command {
	create := &cobra.Command{
		Use:   "create [resource]",
		Short: "创建资源",
	}
	create.AddCommand(c.token.Create())

	get := &cobra.Command{
		Use:   "get [resource]",
		Short: "获取资源信息",
	}
	get.AddCommand(c.script.Get(), c.value.Get(), c.cookie.Get(), c.token.Get())

	edit := &cobra.Command{
		Use:   "edit [resource]",
		Short: "编辑资源信息",
	}
	edit.AddCommand(c.script.Edit())

	del := &cobra.Command{
		Use:   "delete [resource]",
		Short: "删除资源信息",
	}
	del.AddCommand(c.script.Delete(), c.token.Delete(), c.value.Delete(), c.cookie.Delete())

	importer := &cobra.Command{
		Use:   "import [resource]",
		Short: "导入资源信息",
	}
	importer.AddCommand(c.cookie.Import(), c.value.Import())

	cmd := []*cobra.Command{create, get, edit, del, importer}
	cmd = append(cmd, c.script.Command()...)
	cmd = append(cmd, c.value.Command()...)
	cmd = append(cmd, c.cookie.Command()...)
	return cmd
}
