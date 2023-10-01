package scripts

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/model/entity/script_entity"
)

type Script struct {
	ID           string                    `json:"id" yaml:"id,omitempty"`
	Name         string                    `json:"name" yaml:"name"`
	Code         string                    `json:"code,omitempty" yaml:"code,omitempty"`
	Metadata     script_entity.Metadata    `json:"metadata" yaml:"metadata"`
	SelfMetadata script_entity.Metadata    `json:"self_metadata" yaml:"selfMetadata"`
	Status       script_entity.Status      `json:"status" yaml:"status"`
	State        script_entity.ScriptState `json:"state" yaml:"state"`
	Createtime   int64                     `json:"createtime" yaml:"createtime"`
	Updatetime   int64                     `json:"updatetime" yaml:"updatetime"`
}

func (s *Script) Entity() *script_entity.Script {
	return &script_entity.Script{
		ID:           s.ID,
		Name:         s.Name,
		Code:         s.Code,
		Metadata:     s.Metadata,
		SelfMetadata: s.SelfMetadata,
		Status:       s.Status,
		State:        s.State,
		Createtime:   s.Createtime,
		Updatetime:   s.Updatetime,
	}
}

// ListRequest 脚本列表
type ListRequest struct {
	mux.Meta `path:"/scripts" method:"GET"`
	ScriptID string `form:"scriptId"`
}

// ListResponse 脚本列表
type ListResponse struct {
	List []*Script `json:"list"`
}

// InstallRequest 创建脚本
type InstallRequest struct {
	mux.Meta `path:"/scripts" method:"POST"`
	Code     string `form:"code"`
}

// InstallResponse 创建脚本
type InstallResponse struct {
	Scripts []*Script `json:"scripts"`
}

// GetRequest 获取脚本
type GetRequest struct {
	mux.Meta `path:"/scripts/:scriptId" method:"GET"`
	ScriptID string `uri:"scriptId"`
}

// GetResponse 获取脚本
type GetResponse struct {
	Script *Script `json:"script"`
}

// UpdateRequest 更新脚本
type UpdateRequest struct {
	mux.Meta `path:"/scripts/:scriptId" method:"PUT"`
	ScriptID string  `uri:"scriptId"`
	Script   *Script `form:"script"`
}

// UpdateResponse 更新脚本
type UpdateResponse struct {
}

// DeleteRequest 删除脚本
type DeleteRequest struct {
	mux.Meta `path:"/scripts/:scriptId" method:"DELETE"`
	ScriptID string `uri:"scriptId"`
}

// DeleteResponse 删除脚本
type DeleteResponse struct {
}

type Storage struct {
	Name         string   `json:"name"`
	LinkScriptID []string `json:"link_script_id"`
}

// StorageListRequest 值储存空间列表
type StorageListRequest struct {
	mux.Meta `path:"/storages" method:"GET"`
}

// StorageListResponse 值储存空间列表
type StorageListResponse struct {
	List []*Storage `json:"list"`
}

// RunRequest 运行脚本
type RunRequest struct {
	mux.Meta `path:"/scripts/:scriptId/run" method:"POST"`
	ScriptID string `uri:"scriptId"`
}

// RunResponse 运行脚本
type RunResponse struct {
}

// StopRequest 停止脚本
type StopRequest struct {
	mux.Meta `path:"/scripts/:scriptId/stop" method:"POST"`
	ScriptID string `uri:"scriptId"`
}

// StopResponse 停止脚本
type StopResponse struct {
}

// WatchRequest 监听脚本
type WatchRequest struct {
	mux.Meta `path:"/scripts/:scriptId/watch" method:"GET"`
	ScriptID string `uri:"scriptId"`
}

// WatchResponse 监听脚本
type WatchResponse struct {
}
