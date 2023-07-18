package script_entity

import (
	"time"

	"github.com/scriptscat/cloudcat/pkg/scriptcat"
)

type Runtime string

type ScriptState string

const (
	ScriptStateEnable  ScriptState = "enable"
	ScriptStateDisable ScriptState = "disable"

	RuntimeScriptCat Runtime = "scriptcat"
)

type Metadata map[string][]string

type Status map[string]string

const (
	RunState         = "runState"
	RunStateRunning  = "running"
	RunStateComplete = "complete"

	ErrorMsg = "errorMsg"
)

type Script struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Code         string      `json:"code"`
	Runtime      Runtime     `json:"runtime"`
	Metadata     Metadata    `json:"metadata"`
	SelfMetadata Metadata    `json:"self_metadata"`
	Status       Status      `json:"status"`
	State        ScriptState `json:"state"`
	CreatedAt    time.Time   `json:"created_time"`
	UpdatedAt    time.Time   `json:"updated_time"`
}

func (s *Script) Create(script *scriptcat.Script) error {
	s.ID = script.ID
	s.Name = script.Metadata["name"][0]
	s.Code = script.Code
	s.Runtime = RuntimeScriptCat
	s.Metadata = Metadata(script.Metadata)
	s.Status = nil
	s.State = ScriptStateEnable
	s.CreatedAt = time.Now()
	return nil
}

func (s *Script) Update(script *scriptcat.Script) error {
	s.Name = script.Metadata["name"][0]
	s.Code = script.Code
	s.Runtime = RuntimeScriptCat
	s.Metadata = Metadata(script.Metadata)
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Script) Scriptcat() *scriptcat.Script {
	return &scriptcat.Script{
		ID:       s.ID,
		Code:     s.Code,
		Metadata: scriptcat.Metadata(s.Metadata),
	}
}

func (s *Script) StorageName() string {
	storageNames, ok := s.Metadata["storageName"]
	if !ok {
		storageNames = []string{s.ID}
	}
	return storageNames[0]
}

func (s Status) GetRunStatus() string {
	if s == nil {
		return ""
	}
	status, ok := s["runStatus"]
	if !ok {
		return RunStateComplete
	}
	return status
}

func (s Status) SetRunStatus(status string) {
	if s == nil {
		return
	}
	s["runStatus"] = status
}
