package command

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/internal/model/entity/cookie_entity"
	"github.com/scriptscat/cloudcat/pkg/cloudcat_api"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
)

type Cookie struct {
}

func NewCookie() *Cookie {
	return &Cookie{}
}

func (c *Cookie) Command() []*cobra.Command {

	return []*cobra.Command{}
}

func (c *Cookie) Get() *cobra.Command {
	ret := &cobra.Command{
		Use:   "cookie [storageName] [host]",
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
					if strings.Contains(v.Host, args[1]) {
						for _, v := range v.Cookies {
							r.WriteLine(v)
						}
					}
				}
				r.Render()
				return nil
			}
			utils.DealTable([]string{
				"HOST",
			}, list.List, func(i interface{}) []string {
				v := i.(*scripts.Cookie)
				return []string{
					v.Host,
				}
			}).Render()
			return nil
		},
	}
	return ret
}

func (c *Cookie) Delete() *cobra.Command {
	ret := &cobra.Command{
		Use:   "cookie [storageName] [host]",
		Short: "删除cookie信息",
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
				for _, v := range list.List {
					if strings.Contains(v.Host, args[1]) {
						if _, err := cli.DeleteCookie(context.Background(), &scripts.DeleteCookieRequest{
							StorageName: storageName,
							Host:        v.Host,
						}); err != nil {
							return err
						}
					}
				}
				return nil
			}
			for _, v := range list.List {
				if _, err := cli.DeleteCookie(context.Background(), &scripts.DeleteCookieRequest{
					StorageName: storageName,
					Host:        v.Host,
				}); err != nil {
					return err
				}
			}
			return nil
		},
	}

	return ret
}

func (c *Cookie) Import() *cobra.Command {
	return &cobra.Command{
		Use:   "cookie [storageName] [file]",
		Short: "导入cookie信息",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := cloudcat_api.NewCookie(cloudcat_api.DefaultClient())
			data, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}
			storageName := args[0]
			// 获取值列表
			m := make([]*cookie_entity.HttpCookie, 0)
			if err := json.Unmarshal(data, &m); err != nil {
				return err
			}
			for _, v := range m {
				if v.Expires.IsZero() && v.ExpirationDate > 0 {
					v.Expires = time.Unix(v.ExpirationDate, 0)
				}
			}
			if _, err := cli.SetCookie(context.Background(), &scripts.SetCookieRequest{
				StorageName: storageName,
				Cookies:     m,
			}); err != nil {
				return err
			}
			return nil
		},
	}
}
