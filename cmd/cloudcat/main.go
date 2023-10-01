package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/scriptscat/cloudcat/cmd/cloudcat/server"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
)

var configFile = "~/.cloudcat/config.yaml"

func init() {
	// 判断是否为windows
	if runtime.GOOS == "windows" {
		configFile = "./cloudcat/config.yaml"
	}
}

func main() {
	var config string
	rootCmd := &cobra.Command{
		Use:   "cloudcat",
		Short: "cloudcat service.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 转为绝对路径
			var err error
			config, err = utils.Abs(config)
			if err != nil {
				return fmt.Errorf("convert to absolute path: %w", err)
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", configFile, "config file")

	serverCmd := server.NewServer()
	rootCmd.AddCommand(serverCmd.Command(&config)...)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute err: %v", err)
	}
}
