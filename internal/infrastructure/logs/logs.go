package logs

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogs(debug bool) {
	logrus.SetReportCaller(true)
	var w io.Writer = &lumberjack.Logger{
		Filename:   "./runtime/logs/runtime.log",
		MaxSize:    2,
		MaxBackups: 10,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   false,
	}
	f := &logrus.JSONFormatter{}
	logrus.SetFormatter(f)
	if debug {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		w = io.MultiWriter(w, os.Stdout)
	}
	logrus.SetOutput(w)
	logrus.AddHook(NewErrorFile(io.MultiWriter(w, &lumberjack.Logger{
		Filename:   "./runtime/errors/errors.log",
		MaxSize:    2,
		MaxBackups: 30,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   false,
	}), f))

	logrus.Infof("init logs")
}
