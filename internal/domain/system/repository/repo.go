package repository

import (
	"context"

	"github.com/scriptscat/cloudcat/pkg/kvdb"
)

type Repo interface {
	GetScriptCatInfo() (*ScriptCatInfo, error)
}

type repo struct {
	kv kvdb.KvDb
}

func NewRepo(kv kvdb.KvDb) Repo {
	return &repo{
		kv: kv,
	}
}

func (r *repo) GetScriptCatInfo() (*ScriptCatInfo, error) {
	v, err := r.kv.Get(context.Background(), "cloudcat:scriptcat:info:version")
	if err != nil {
		return nil, err
	}
	n, err := r.kv.Get(context.Background(), "cloudcat:scriptcat:info:notice")
	if err != nil {
		return nil, err
	}
	ret := &ScriptCatInfo{
		Version: v,
		Notice:  n,
	}
	return ret, nil
}
