package scripts

import (
	"time"

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
	CreatedTime  time.Time                 `json:"created_time" yaml:"createdTime"`
	UpdatedTime  time.Time                 `json:"updated_time" yaml:"updatedTime"`
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
