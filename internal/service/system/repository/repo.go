package repository

import (
	"context"
	"encoding/json"

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
	s, err := r.kv.Get(context.Background(), "cloudcat:scriptcat:info")
	if err != nil {
		return nil, err
	}
	ret := &ScriptCatInfo{}
	if err := json.Unmarshal([]byte(s), ret); err != nil {
		return nil, err
	}
	return ret, nil
}
