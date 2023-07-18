package scripts_ctr

import (
	"context"

	api "github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/internal/service/scripts_svc"
)

type Value struct {
}

func NewValue() *Value {
	return &Value{}
}

// ValueList 脚本值列表
func (v *Value) ValueList(ctx context.Context, req *api.ValueListRequest) (*api.ValueListResponse, error) {
	return scripts_svc.Value().ValueList(ctx, req)
}

// StorageList 值储存空间列表
func (v *Value) StorageList(ctx context.Context, req *api.StorageListRequest) (*api.StorageListResponse, error) {
	return scripts_svc.Value().StorageList(ctx, req)
}
