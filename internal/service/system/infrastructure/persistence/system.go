package persistence

import (
	"context"

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
	v, err := r.kv.Get(context.Background(), "cloudcat:scriptcat:info:version")
	if err != nil {
		return nil, err
	}
	n, err := r.kv.Get(context.Background(), "cloudcat:scriptcat:info:notice")
	if err != nil {
		return nil, err
	}
	ret := &repository.ScriptCatInfo{
		Version: v,
		Notice:  n,
	}
	return ret, nil
}
