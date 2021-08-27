package system

import (
	"github.com/scriptscat/cloudcat/internal/service/system/repository"
)

type System struct {
	repo repository.Repo
}

func NewSystem(repo repository.Repo) *System {
	return &System{repo: repo}
}

func (s *System) ScriptCatInfo() (*repository.ScriptCatInfo, error) {
	return s.repo.GetScriptCatInfo()
}
