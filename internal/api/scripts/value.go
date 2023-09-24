package scripts

import (
	"github.com/codfrm/cago/server/mux"
)

type Value struct {
	StorageName string `json:"storage_name"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Createtime  int64  `json:"createtime"`
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
