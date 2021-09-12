package service

import (
	"github.com/scriptscat/cloudcat/internal/domain/system/repository"
)

type System interface {
	ScriptCatInfo() (*repository.ScriptCatInfo, error)
}

const (
	SYSTEM_SITE_NAME = "system_site_name"
)

type system struct {
	repo repository.Repo
}

func NewSystem(repo repository.Repo) System {
	return &system{repo: repo}
}

func (s *system) ScriptCatInfo() (*repository.ScriptCatInfo, error) {
	return s.repo.GetScriptCatInfo()
}
