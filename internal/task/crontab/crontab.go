package crontab

import (
	"github.com/codfrm/cago/server/cron"
	"github.com/scriptscat/cloudcat/internal/task/crontab/handler"
)

type Cron interface {
	Crontab(c cron.Crontab) error
}

// Crontab 定时任务
func Crontab(cron cron.Crontab) error {
	crontab := []Cron{&handler.Script{}}
	for _, v := range crontab {
		if err := v.Crontab(cron); err != nil {
			return err
		}
	}
	return nil
}
