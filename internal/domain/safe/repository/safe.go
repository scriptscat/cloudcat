package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
)

type Safe interface {
	GetLastOpTime(user, op string) (int64, error)
	GetPeriodOpCnt(user, op string) (int64, error)
	SetLastOpTime(user, op string, t int64, expired time.Duration) error
	DelLastOpTime(user, op string) error
}

type rate struct {
	kv kvdb.KvDb
}

func NewRate(kv kvdb.KvDb) Safe {
	return &rate{
		kv: kv,
	}
}

func (r *rate) GetLastOpTime(user, op string) (int64, error) {
	ret, err := r.kv.Get(context.Background(), r.key(user, op))
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(ret, 10, 64)
}

func (r *rate) SetLastOpTime(user, op string, t int64, expired time.Duration) error {
	err := r.kv.Set(context.Background(), r.key(user, op), strconv.FormatInt(t, 10), time.Hour)
	if err != nil {
		return err
	}
	k := r.key(user, op) + ":period"
	if err := r.kv.IncrBy(context.Background(), k, 1); err != nil {
		return err
	}
	if err := r.kv.Expire(context.Background(), k, expired); err != nil {
		return err
	}
	return nil
}

func (r *rate) DelLastOpTime(user, op string) error {
	err := r.kv.Set(context.Background(), r.key(user, op), "", time.Hour)
	if err != nil {
		return err
	}
	k := r.key(user, op) + ":period"
	if err := r.kv.IncrBy(context.Background(), k, -1); err != nil {
		return err
	}
	return nil
}

func (r *rate) GetPeriodOpCnt(user, op string) (int64, error) {
	k := r.key(user, op) + ":period"
	ret, err := r.kv.Get(context.Background(), k)
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return strconv.ParseInt(ret, 10, 64)
}

func (r *rate) key(user, op string) string {
	return fmt.Sprintf("safe:rate:" + user + ":" + op)
}
