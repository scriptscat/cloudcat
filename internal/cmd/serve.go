package cmd

import (
	"fmt"

	"github.com/scriptscat/cloudcat/internal/app"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/spf13/cobra"
)

type serveCmd struct {
	config string
}

func NewServeCmd() *serveCmd {
	return &serveCmd{}
}

func (s *serveCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "serve [flag]",
		Short: "运行脚本猫服务",
		RunE:  s.serve,
	}
	ret.Flags().StringVarP(&s.config, "config", "c", "config.yaml", "配置文件")

	return []*cobra.Command{ret}
}

func (s *serveCmd) serve(cmd *cobra.Command, args []string) error {
	cfg, err := config.Init(s.config)
	if err != nil {
		return err
	}

	if err := app.Run(cfg); err != nil {
		return fmt.Errorf("app start err: %v", err)
	}
	return nil
}
