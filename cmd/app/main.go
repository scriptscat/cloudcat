package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "cloudcat",
	}

	execCmd := newExecCmd()
	serveCmd := newServeCmd()

	rootCmd.AddCommand(execCmd.Commands()...)
	rootCmd.AddCommand(serveCmd.Commands()...)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
