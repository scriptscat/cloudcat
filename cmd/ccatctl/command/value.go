package command

import (
	"context"
	"strings"

	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/pkg/cloudcat_api"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
)

type Value struct {
	cli    *mux.Client
	config *string
}

func NewValue(config *string) *Value {
	return &Value{
		cli:    mux.NewClient("http://127.0.0.1:8080/api/v1"),
		config: config,
	}
}

func (s *Value) Command() []*cobra.Command {

	return []*cobra.Command{}
}

func (s *Value) Get() *cobra.Command {
	ret := &cobra.Command{
		Use:   "value [storageName/scriptId]",
		Short: "获取值信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := cloudcat_api.NewValue(s.cli)
			storageName := ""
			if len(args) > 0 {
				storageName = args[0]
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
			}
			list, err := cli.StorageList(context.Background(), &scripts.StorageListRequest{})
			if err != nil {
				return err
			}
			utils.DealTable([]string{
				"ID", "LINK",
			}, list.List, func(i interface{}) []string {
				v := i.(*scripts.Storage)
				name := v.Name
				if len(name) > 7 {
					name = name[:7]
				}
				link := make([]string, 0)
				for k := range v.LinkScriptID {
					link = append(link, k)
				}
				return []string{name,
					v.Name, strings.Join(link, ","),
				}
			}).Render()
			return nil
		},
	}
	return ret
}
