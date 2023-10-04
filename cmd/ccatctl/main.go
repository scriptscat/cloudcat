package main

import (
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/scriptscat/cloudcat/configs"

	"github.com/scriptscat/cloudcat/cmd/ccatctl/command"
	"github.com/scriptscat/cloudcat/pkg/cloudcat_api"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configFile = "~/.cloudcat/cloudcat.yaml"

func init() {
	// 判断是否为windows
	if runtime.GOOS == "windows" {
		configFile = "./cloudcat/cloudcat.yaml"
	}
}

func main() {
	config := ""
	rootCmd := &cobra.Command{
		Use:     "ccatctl",
		Short:   "ccatctl controls the cloudcat service.",
		Version: configs.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.Abs(config)
			if err != nil {
				return err
			}
			_, err = os.Stat(config)
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				// 从环境变量中获取
				var ok bool
				config, ok = os.LookupEnv("CCATCONFIG")
				if !ok {
					return errors.New("config file is not exist")
				}
			}
			configData, err := os.ReadFile(config)
			if err != nil {
				return err
			}
			cfg := &cloudcat_api.Config{}
			if err := yaml.Unmarshal(configData, cfg); err != nil {
				return err
			}
			cli := cloudcat_api.NewClient(cfg)
			cloudcat_api.SetDefaultClient(cli)
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", configFile, "config file")
	basic := command.NewBasic(config)
	rootCmd.AddCommand(basic.Command()...)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute err: %v", err)
	}
}
