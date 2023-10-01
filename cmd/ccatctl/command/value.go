package command

import (
	"context"
	"encoding/json"
	"os"

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
				value, _ := json.Marshal(v.Value.Get())
				return []string{
					v.Key, string(value),
				}
			}).Render()
			return nil
		},
	}
	return ret
}

func (s *Value) Import() *cobra.Command {
	return &cobra.Command{
		Use:   "value [storageName] [file]",
		Short: "导入值信息",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := cloudcat_api.NewValue(cloudcat_api.DefaultClient())
			data, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}
			storageName := args[0]
			// 获取值列表
			m := make([]*scripts.Value, 0)
			if err := json.Unmarshal(data, &m); err != nil {
				return err
			}
			if _, err := cli.SetValue(context.Background(), &scripts.SetValueRequest{
				StorageName: storageName,
				Values:      m,
			}); err != nil {
				return err
			}
			return nil
		},
	}
}

func (s *Value) Delete() *cobra.Command {
	return &cobra.Command{
		Use:   "value [storageName] [key]",
		Short: "删除值信息",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := cloudcat_api.NewValue(cloudcat_api.DefaultClient())
			storageName := args[0]
			key := args[1]
			// 获取值列表
			if _, err := cli.DeleteValue(context.Background(), &scripts.DeleteValueRequest{
				StorageName: storageName,
				Key:         key,
			}); err != nil {
				return err
			}
			return nil
		},
	}
}
