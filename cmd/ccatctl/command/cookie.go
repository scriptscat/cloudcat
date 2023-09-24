package command

import (
	"context"
	"github.com/scriptscat/cloudcat/pkg/scriptcat/cookie"
	"strings"

	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/pkg/cloudcat_api"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
)

type Cookie struct {
	cli    *mux.Client
	config *string
}

func NewCookie(config *string) *Cookie {
	return &Cookie{
		cli:    mux.NewClient("http://127.0.0.1:8080/api/v1"),
		config: config,
	}
}

func (s *Cookie) Command() []*cobra.Command {

	return []*cobra.Command{}
}

func (s *Cookie) Get() *cobra.Command {
	ret := &cobra.Command{
		Use:   "cookie [storageName] [url]",
		Short: "获取cookie信息",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := cloudcat_api.NewCookie(s.cli)
			storageName := args[0]
			// 获取值列表
			list, err := cli.CookieList(context.Background(), &scripts.CookieListRequest{
				StorageName: storageName,
			})
			if err != nil {
				return err
			}
			if len(args) > 1 {
				r := utils.DealTable([]string{
					"NAME", "VALUE", "DOMAIN", "PATH", "EXPIRES", "HTTPONLY", "SECURE",
				}, nil, func(i interface{}) []string {
					v := i.(*cookie.Cookie)
					return []string{
						v.Name, v.Value, v.Domain, v.Path,
						v.Expires.Format("2006-01-02 15:04:05"),
						utils.BoolToString(v.HttpOnly), utils.BoolToString(v.Secure),
					}
				})
				for _, v := range list.List {
					if strings.Contains(v.Url, args[1]) {
						for _, v := range v.Cookies {
							r.WriteLine(v)
						}
					}
				}
				r.Render()
				return nil
			}
			utils.DealTable([]string{
				"URL",
			}, list.List, func(i interface{}) []string {
				v := i.(*scripts.Cookie)
				return []string{
					v.Url,
				}
			}).Render()
			return nil
		},
	}
	return ret
}
