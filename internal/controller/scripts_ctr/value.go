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
