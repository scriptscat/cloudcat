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
	Createtime   int64       `json:"createtime"`
	Updatetime   int64       `json:"updatetime"`
}

func (s *Script) Create(script *scriptcat.Script) error {
	s.ID = script.ID
	s.Name = script.Metadata["name"][0]
	s.Code = script.Code
	s.Runtime = RuntimeScriptCat
	s.Metadata = Metadata(script.Metadata)
	s.Status = nil
	s.State = ScriptStateEnable
	s.Createtime = time.Now().Unix()
	return nil
}

func (s *Script) Update(script *scriptcat.Script) error {
	s.Name = script.Metadata["name"][0]
	s.Code = script.Code
	s.Runtime = RuntimeScriptCat
	s.Metadata = Metadata(script.Metadata)
	s.Updatetime = time.Now().Unix()
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
	return StorageName(s.ID, s.Metadata)
}

func StorageName(id string, m Metadata) string {
	storageNames, ok := m["storageName"]
	if !ok {
		storageNames = []string{id}
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
