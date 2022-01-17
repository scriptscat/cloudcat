package logs

import (
	"io"

	"github.com/sirupsen/logrus"
)

type errorFile struct {
	w   io.Writer
	log *logrus.Logger
}

func NewErrorFile(w io.Writer, formatter logrus.Formatter) logrus.Hook {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetOutput(w)
	log.SetFormatter(formatter)
	return &errorFile{
		w:   w,
		log: log,
	}
}

func (e *errorFile) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel}
}

func (e *errorFile) Fire(entry *logrus.Entry) error {
	entry.Logger = e.log
	return nil
}
