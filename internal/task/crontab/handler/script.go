package handler

import (
	"github.com/codfrm/cago/server/cron"
)

type Script struct {
}

func (s *Script) Crontab(c cron.Crontab) error {
	return nil
}
