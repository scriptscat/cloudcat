package service

import (
	"time"

	"github.com/scriptscat/cloudcat/internal/domain/safe/dto"
	"github.com/scriptscat/cloudcat/internal/domain/safe/errs"
	"github.com/scriptscat/cloudcat/internal/domain/safe/repository"
)

type Safe interface {
	Rate(userinfo *dto.SafeUserinfo, rule *dto.SafeRule, f func() error) error
	Limit(userinfo *dto.SafeUserinfo, rule *dto.SafeRule, f func() error) error
}

type rate struct {
	repo repository.Safe
}

func NewRate(repo repository.Safe) Safe {
	return &rate{repo: repo}
}

// Rate 不管成功与否都计算一次
func (r *rate) Rate(userinfo *dto.SafeUserinfo, rule *dto.SafeRule, f func() error) error {
	t, err := r.repo.GetLastOpTime(userinfo.Userinfo(), rule.Name)
	if err != nil {
		return err
	}
	if t > time.Now().Unix()-rule.Interval {
		return errs.NewOperationTimeToShort(rule)
	}
	c, err := r.repo.GetPeriodOpCnt(userinfo.Userinfo(), rule.Name)
	if err != nil {
		return err
	}
	if rule.PeriodCnt > 0 && c > rule.PeriodCnt {
		return errs.NewOperationMax(rule)
	}
	if err := r.repo.SetLastOpTime(userinfo.Userinfo(), rule.Name, time.Now().Unix(), rule.Period); err != nil {
		return err
	}
	return f()
}

// Limit 成功才会计算一次
func (r *rate) Limit(userinfo *dto.SafeUserinfo, rule *dto.SafeRule, f func() error) error {
	t, err := r.repo.GetLastOpTime(userinfo.Userinfo(), rule.Name)
	if err != nil {
		return err
	}
	if t > time.Now().Unix()-rule.Interval {
		return errs.NewOperationTimeToShort(rule)
	}
	c, err := r.repo.GetPeriodOpCnt(userinfo.Userinfo(), rule.Name)
	if err != nil {
		return err
	}
	if rule.PeriodCnt > 0 && c > rule.PeriodCnt {
		return errs.NewOperationLimit(rule)
	}
	if err := r.repo.SetLastOpTime(userinfo.Userinfo(), rule.Name, time.Now().Unix(), rule.Period); err != nil {
		return err
	}
	if err := f(); err != nil {
		_ = r.repo.DelLastOpTime(userinfo.Userinfo(), rule.Name)
		return err
	}
	return nil
}
