package dto

import (
	"strconv"
	"time"
)

type SafeRule struct {
	Name        string
	Description string
	Interval    int64
	PeriodCnt   int64
	Period      time.Duration
}

type SafeUserinfo struct {
	Identifier string
	Uid        int64
	IP         string
	UserAgent  string
}

func (r *SafeUserinfo) Userinfo() string {
	if r.Identifier != "" {
		return r.Identifier
	}
	if r.Uid != 0 {
		return strconv.FormatInt(r.Uid, 10)
	}
	if r.UserAgent != "" {
		return r.IP + "-" + r.UserAgent
	}
	return r.IP
}
