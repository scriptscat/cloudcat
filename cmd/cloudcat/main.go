package main

import (
	cmd2 "github.com/scriptscat/cloudcat/cmd/cloudcat/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "cloudcat",
	}

	execCmd := cmd2.NewExecCmd()
	serveCmd := cmd2.NewServeCmd()
	manageCmd := cmd2.NewManageCmd()

	rootCmd.AddCommand(execCmd.Commands()...)
	rootCmd.AddCommand(serveCmd.Commands()...)
	rootCmd.AddCommand(manageCmd.Commands()...)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
