package persistence

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/cloudcat/internal/pkg/kvdb"
	"github.com/scriptscat/cloudcat/internal/service/system/domain/repository"
)

type repo struct {
	kv kvdb.KvDb
}

func NewRepo(kv kvdb.KvDb) repository.System {
	return &repo{
		kv: kv,
	}
}

func (r *repo) GetScriptCatInfo() (*repository.ScriptCatInfo, error) {
	ret := &repository.ScriptCatInfo{}
	var err error
	ret.Version, err = r.kv.Get(context.Background(), "cloudcat:scriptcat:info:version")
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
	}
	ret.Notice, err = r.kv.Get(context.Background(), "cloudcat:scriptcat:info:notice")
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
	}
	return ret, nil
}
