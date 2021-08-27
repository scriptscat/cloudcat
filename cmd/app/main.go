package main

import (
	"github.com/scriptscat/cloudcat/internal/app"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Init("config.yaml")
	if err != nil {
		logrus.Fatalln(err)
	}

	if err := app.Run(cfg); err != nil {
		logrus.Fatalf("app start err: %v", err)
	}
}
