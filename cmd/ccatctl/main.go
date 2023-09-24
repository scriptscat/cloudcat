package main

import (
	"log"

	"github.com/scriptscat/cloudcat/cmd/ccatctl/command"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ccatctl",
		Short: "ccatctl controls the cloudcat service.",
	}
	config := rootCmd.PersistentFlags().StringP("config", "c", "./configs/config.yaml", "config file")
	cmd := command.NewBasic(config)
	rootCmd.AddCommand(cmd.Command()...)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute err: %v", err)
	}
}
