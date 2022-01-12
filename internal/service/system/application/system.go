package application

import (
	"github.com/scriptscat/cloudcat/internal/service/system/domain/repository"
)

type System interface {
	ScriptCatInfo() (*repository.ScriptCatInfo, error)
}

const (
	SYSTEM_SITE_NAME = "system_site_name"
)

type system struct {
	repo repository.System
}

func NewSystem(repo repository.System) System {
	return &system{repo: repo}
}

func (s *system) ScriptCatInfo() (*repository.ScriptCatInfo, error) {
	return s.repo.GetScriptCatInfo()
}
