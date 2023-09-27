package command

import (
	"context"

	"github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/pkg/cloudcat_api"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
)

type Value struct {
}

func NewValue() *Value {
	return &Value{}
}

func (s *Value) Command() []*cobra.Command {

	return []*cobra.Command{}
}

func (s *Value) Get() *cobra.Command {
	ret := &cobra.Command{
		Use:   "value [storageName]",
		Short: "获取值信息",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := cloudcat_api.NewValue(cloudcat_api.DefaultClient())
			storageName := args[0]
			// 获取值列表
			list, err := cli.ValueList(context.Background(), &scripts.ValueListRequest{
				StorageName: storageName,
			})
			if err != nil {
				return err
			}
			utils.DealTable([]string{
				"KEY", "VALUE",
			}, list.List, func(i interface{}) []string {
				v := i.(*scripts.Value)
				return []string{
					v.Key, v.Value,
				}
			}).Render()
			return nil
		},
	}
	return ret
}
