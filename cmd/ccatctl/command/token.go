package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/scriptscat/cloudcat/internal/api/auth"
	"github.com/scriptscat/cloudcat/pkg/cloudcat_api"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Token struct {
	out string
}

func NewToken() *Token {
	return &Token{}
}

func (s *Token) Command() []*cobra.Command {
	return []*cobra.Command{}
}

func (s *Token) Get() *cobra.Command {
	ret := &cobra.Command{
		Use:   "token [tokenId]",
		Short: "获取脚本信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := cloudcat_api.NewToken(cloudcat_api.DefaultClient())
			tokenId := ""
			if len(args) > 0 {
				tokenId = args[0]
			}
			list, err := cli.List(context.Background(), &auth.TokenListRequest{
				TokenID: tokenId,
			})
			if err != nil {
				return err
			}
			if s.out == "yaml" {
				for _, v := range list.List {
					data, err := yaml.Marshal(v)
					if err != nil {
						return err
					}
					_, err = os.Stdout.Write(data)
					if err != nil {
						return err
					}
				}
				return nil
			}
			utils.DealTable([]string{
				"ID", "CREATED_AT",
			}, list.List, func(i interface{}) []string {
				v := i.(*auth.Token)
				return []string{
					v.ID,
					time.Unix(v.Createtime, 0).Format("2006-01-02 15:04:05"),
				}
			}).Render()
			return nil
		},
	}
	ret.Flags().StringVarP(&s.out, "out", "o", "table", "输出格式: yaml, table")
	return ret
}

func (s *Token) Delete() *cobra.Command {
	ret := &cobra.Command{
		Use:   "token [tokenId]",
		Short: "删除脚本",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenId := args[0]
			cli := cloudcat_api.NewToken(cloudcat_api.DefaultClient())
			_, err := cli.Delete(context.Background(), &auth.TokenDeleteRequest{
				TokenID: tokenId,
			})
			if err != nil {
				return err
			}
			return nil
		},
	}
	return ret
}

func (s *Token) Create() *cobra.Command {
	ret := &cobra.Command{
		Use:   "token [tokenId]",
		Short: "创建脚本",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenId := args[0]
			cli := cloudcat_api.NewToken(cloudcat_api.DefaultClient())
			resp, err := cli.Create(context.Background(), &auth.TokenCreateRequest{
				TokenID: tokenId,
			})
			if err != nil {
				return err
			}
			data, err := yaml.Marshal(resp.Token)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		},
	}
	return ret
}
