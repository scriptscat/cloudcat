package main

import (
	"github.com/scriptscat/cloudcat/internal/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "cloudcat",
	}

	execCmd := cmd.NewExecCmd()
	serveCmd := cmd.NewServeCmd()
	manageCmd := cmd.NewManageCmd()

	rootCmd.AddCommand(execCmd.Commands()...)
	rootCmd.AddCommand(serveCmd.Commands()...)
	rootCmd.AddCommand(manageCmd.Commands()...)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
