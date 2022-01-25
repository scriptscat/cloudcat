package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/cloudcat/internal/pkg/kvdb"
	"github.com/scriptscat/cloudcat/pkg/utils"
)

type Safe interface {
	GetLastOpTime(user, op string) (int64, error)
	GetPeriodOpCnt(user, op string) (int64, error)
	SetLastOpTime(user, op string, t int64, expired time.Duration) error
	DelLastOpTime(user, op string) error
}

type safe struct {
	kv kvdb.KvDb
}

func NewSafe(kv kvdb.KvDb) Safe {
	return &safe{
		kv: kv,
	}
}

func (s *safe) GetLastOpTime(user, op string) (int64, error) {
	ret, err := s.kv.Get(context.Background(), s.key(user, op))
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return utils.StringToInt64(ret), nil
}

func (s *safe) SetLastOpTime(user, op string, t int64, expired time.Duration) error {
	err := s.kv.Set(context.Background(), s.key(user, op), strconv.FormatInt(t, 10), time.Hour)
	if err != nil {
		return err
	}
	k := s.key(user, op) + ":period"
	if _, err := s.kv.IncrBy(context.Background(), k, 1); err != nil {
		return err
	}
	if err := s.kv.Expire(context.Background(), k, expired); err != nil {
		return err
	}
	return nil
}

func (s *safe) DelLastOpTime(user, op string) error {
	err := s.kv.Set(context.Background(), s.key(user, op), "", time.Hour)
	if err != nil {
		return err
	}
	k := s.key(user, op) + ":period"
	if _, err := s.kv.IncrBy(context.Background(), k, -1); err != nil {
		return err
	}
	return nil
}

func (s *safe) GetPeriodOpCnt(user, op string) (int64, error) {
	k := s.key(user, op) + ":period"
	ret, err := s.kv.Get(context.Background(), k)
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return utils.StringToInt64(ret), nil
}

func (s *safe) key(user, op string) string {
	return fmt.Sprintf("safe:safe:" + user + ":" + op)
}
