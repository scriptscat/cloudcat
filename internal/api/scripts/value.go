package scripts

import (
	"time"

	"github.com/codfrm/cago/server/mux"
)

type Value struct {
	StorageName string    `json:"storage_name"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	CreatedTime time.Time `json:"created_time"`
}

type Storage struct {
	Name         string   `json:"name"`
	LinkScriptID []string `json:"link_script_id"`
}

// StorageListRequest 值储存空间列表
type StorageListRequest struct {
	mux.Meta `path:"/values" method:"GET"`
}

// StorageListResponse 值储存空间列表
type StorageListResponse struct {
	List []*Storage `json:"list"`
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
