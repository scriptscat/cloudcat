package scripts

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/model/entity/value_entity"
)

type Value struct {
	StorageName string                   `json:"storage_name"`
	Key         string                   `json:"key"`
	Value       value_entity.ValueString `json:"value"`
	Createtime  int64                    `json:"createtime"`
}

// ValueListRequest 脚本值列表
type ValueListRequest struct {
	mux.Meta    `path:"/values/:storageName" method:"GET"`
	StorageName string `uri:"storageName"`
}

// ValueListResponse 脚本值列表
type ValueListResponse struct {
	List []*Value `json:"list"`
}

// SetValueRequest 设置脚本值
type SetValueRequest struct {
	mux.Meta    `path:"/values/:storageName" method:"POST"`
	StorageName string   `uri:"storageName"`
	Values      []*Value `form:"values"`
}

// SetValueResponse 设置脚本值
type SetValueResponse struct {
}

// DeleteValueRequest 删除脚本值
type DeleteValueRequest struct {
	mux.Meta    `path:"/values/:storageName/:key" method:"DELETE"`
	StorageName string `uri:"storageName"`
	Key         string `uri:"key"`
}

// DeleteValueResponse 删除脚本值
type DeleteValueResponse struct {
}
