package command

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/scriptscat/cloudcat/internal/model/entity/script_entity"

	"github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/pkg/cloudcat_api"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Script struct {
	file string
	out  string
}

func NewScript() *Script {
	return &Script{}
}

func (s *Script) Command() []*cobra.Command {
	install := &cobra.Command{
		Use:   "install",
		Short: "安装脚本",
		RunE:  s.install,
	}
	install.Flags().StringVarP(&s.file, "file", "f", "", "脚本文件")
	run := &cobra.Command{
		Use:   "run [script]",
		Short: "运行脚本",
		Args:  cobra.ExactArgs(1),
		RunE:  s.run,
	}
	stop := &cobra.Command{
		Use:   "stop [script]",
		Short: "停止脚本",
		Args:  cobra.ExactArgs(1),
		RunE:  s.stop,
	}
	enable := &cobra.Command{
		Use:   "enable [script]",
		Short: "启用脚本",
		Args:  cobra.ExactArgs(1),
		RunE:  s.enable(true),
	}
	disable := &cobra.Command{
		Use:   "disable [script]",
		Short: "禁用脚本",
		Args:  cobra.ExactArgs(1),
		RunE:  s.enable(false),
	}

	return []*cobra.Command{install, run, stop, enable, disable}
}

func (s *Script) enable(enable bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cli := cloudcat_api.NewScript(cloudcat_api.DefaultClient())
		script, err := cli.Get(context.Background(), &scripts.GetRequest{
			ScriptID: args[0],
		})
		if err != nil {
			return err
		}
		if enable {
			script.Script.State = script_entity.ScriptStateEnable
		} else {
			script.Script.State = script_entity.ScriptStateDisable
		}
		if _, err := cli.Update(context.Background(), &scripts.UpdateRequest{
			ScriptID: args[0],
			Script:   script.Script,
		}); err != nil {
			return err
		}
		return nil
	}
}

func (s *Script) Get() *cobra.Command {
	ret := &cobra.Command{
		Use:   "script [scriptId]",
		Short: "获取脚本信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := cloudcat_api.NewScript(cloudcat_api.DefaultClient())
			scriptId := ""
			if len(args) > 0 {
				scriptId = args[0]
			}
			list, err := cli.List(context.Background(), &scripts.ListRequest{
				ScriptID: scriptId,
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
				"ID", "NAME", "STORAGE_NAME", "RUN_AT", "CREATED_AT",
			}, list.List, func(i interface{}) []string {
				v := i.(*scripts.Script)
				sn := script_entity.StorageName(v.ID, v.Metadata)
				if len(sn) > 7 {
					sn = sn[:7]
				}
				runAt := ""
				if v.State == script_entity.ScriptStateDisable {
					runAt = "DISABLE"
				} else {
					if cron, ok := v.Entity().Crontab(); ok {
						runAt = cron
					} else {
						runAt = "BACKGROUND"
					}
				}
				return []string{
					v.ID[:7],
					v.Name,
					sn,
					runAt,
					time.Unix(v.Createtime, 0).Format("2006-01-02 15:04:05"),
				}
			}).Render()
			return nil
		},
	}
	ret.Flags().StringVarP(&s.out, "out", "o", "table", "输出格式: yaml, table")
	return ret
}

const (
	defaultEditor = "vi"
	defaultShell  = "/bin/bash"
)

func (s *Script) Edit() *cobra.Command {
	ret := &cobra.Command{
		Use:   "script [scriptId]",
		Short: "编辑脚本信息",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scriptId := args[0]
			cli := cloudcat_api.NewScript(cloudcat_api.DefaultClient())
			resp, err := cli.Get(context.Background(), &scripts.GetRequest{
				ScriptID: scriptId,
			})
			if err != nil {
				return err
			}
			data, err := yaml.Marshal(resp.Script)
			if err != nil {
				return err
			}
			file, err := os.CreateTemp(os.TempDir(), "")
			if err != nil {
				return err
			}
			defer func() {
				_ = file.Close()
				_ = os.Remove(file.Name())
			}()
			// 联合vi编辑
			_, err = file.Write(data)
			if err != nil {
				return err
			}
			c := exec.Command(defaultShell, "-c", defaultEditor+" "+file.Name())
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			err = c.Run()
			if err != nil {
				return err
			}
			editData, err := os.ReadFile(file.Name())
			if err != nil {
				return err
			}
			if bytes.Equal(editData, data) {
				return nil
			}
			script := &scripts.Script{}
			err = yaml.Unmarshal(editData, script)
			if err != nil {
				return err
			}
			if _, err := cli.Update(context.Background(), &scripts.UpdateRequest{
				ScriptID: script.ID,
				Script:   script,
			}); err != nil {
				return err
			}
			return nil
		},
	}
	return ret
}

func (s *Script) install(cmd *cobra.Command, args []string) error {
	cli := cloudcat_api.NewScript(cloudcat_api.DefaultClient())
	code, err := os.ReadFile(s.file)
	if err != nil {
		return err
	}
	_, err = cli.Install(context.Background(), &scripts.InstallRequest{
		Code: string(code),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Script) run(cmd *cobra.Command, args []string) error {
	cli := cloudcat_api.NewScript(cloudcat_api.DefaultClient())
	if _, err := cli.Run(context.Background(), &scripts.RunRequest{
		ScriptID: args[0],
	}); err != nil {
		return err
	}
	return nil
}

func (s *Script) stop(cmd *cobra.Command, args []string) error {
	cli := cloudcat_api.NewScript(cloudcat_api.DefaultClient())
	if _, err := cli.Stop(context.Background(), &scripts.StopRequest{
		ScriptID: args[0],
	}); err != nil {
		return err
	}
	return nil
}

func (s *Script) Delete() *cobra.Command {
	ret := &cobra.Command{
		Use:   "script [scriptId]",
		Short: "删除脚本",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scriptId := args[0]
			cli := cloudcat_api.NewScript(cloudcat_api.DefaultClient())
			_, err := cli.Delete(context.Background(), &scripts.DeleteRequest{
				ScriptID: scriptId,
			})
			if err != nil {
				return err
			}
			return nil
		},
	}
	return ret
}
