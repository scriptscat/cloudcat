package main

import (
	"log"

	"github.com/scriptscat/cloudcat/cmd/cloudcat/server"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cloudcat",
		Short: "cloudcat service.",
	}
	config := rootCmd.PersistentFlags().StringP("config", "c", "./configs/config.yaml", "config file")

	serverCmd := server.NewServer()
	rootCmd.AddCommand(serverCmd.Command(config)...)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute err: %v", err)
	}
}
