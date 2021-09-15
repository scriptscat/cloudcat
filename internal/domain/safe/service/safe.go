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

type safe struct {
	repo repository.Safe
}

func NewSafe(repo repository.Safe) Safe {
	return &safe{repo: repo}
}

// Rate 不管成功与否都计算一次
func (s *safe) Rate(userinfo *dto.SafeUserinfo, rule *dto.SafeRule, f func() error) error {
	if userinfo.IP != "" {
		return s.rate(&dto.SafeUserinfo{
			IP: userinfo.IP,
		}, rule, func() error {
			return s.rate(userinfo, rule, f)
		})
	} else {
		return s.rate(userinfo, rule, f)
	}
}

func (s *safe) rate(userinfo *dto.SafeUserinfo, rule *dto.SafeRule, f func() error) error {
	t, err := s.repo.GetLastOpTime(userinfo.Userinfo(), rule.Name)
	if err != nil {
		return err
	}

	if t > time.Now().Unix()-rule.Interval {
		return errs.NewOperationTimeToShort(rule)
	}
	c, err := s.repo.GetPeriodOpCnt(userinfo.Userinfo(), rule.Name)
	if err != nil {
		return err
	}
	if rule.PeriodCnt > 0 && c > rule.PeriodCnt {
		return errs.NewOperationMax(rule)
	}
	if err := s.repo.SetLastOpTime(userinfo.Userinfo(), rule.Name, time.Now().Unix(), rule.Period); err != nil {
		return err
	}
	return f()
}

// Limit 成功才会计算一次
func (s *safe) Limit(userinfo *dto.SafeUserinfo, rule *dto.SafeRule, f func() error) error {
	t, err := s.repo.GetLastOpTime(userinfo.Userinfo(), rule.Name)
	if err != nil {
		return err
	}
	if t > time.Now().Unix()-rule.Interval {
		return errs.NewOperationTimeToShort(rule)
	}
	c, err := s.repo.GetPeriodOpCnt(userinfo.Userinfo(), rule.Name)
	if err != nil {
		return err
	}
	if rule.PeriodCnt > 0 && c > rule.PeriodCnt {
		return errs.NewOperationLimit(rule)
	}
	if err := s.repo.SetLastOpTime(userinfo.Userinfo(), rule.Name, time.Now().Unix(), rule.Period); err != nil {
		return err
	}
	if err := f(); err != nil {
		_ = s.repo.DelLastOpTime(userinfo.Userinfo(), rule.Name)
		return err
	}
	return nil
}
