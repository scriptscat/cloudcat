package command

import (
	"context"
	"net/http"
	"strings"

	"github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/pkg/cloudcat_api"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
)

type Cookie struct {
}

func NewCookie() *Cookie {
	return &Cookie{}
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
			cli := cloudcat_api.NewCookie(cloudcat_api.DefaultClient())
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
					v := i.(*http.Cookie)
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
