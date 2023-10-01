package scripts_svc

import (
	"context"
	"time"

	"github.com/scriptscat/cloudcat/internal/model/entity/value_entity"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/cloudcat/internal/pkg/code"
	"github.com/scriptscat/cloudcat/internal/repository/script_repo"

	api "github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/internal/repository/value_repo"
)

type ValueSvc interface {
	// ValueList 脚本值列表
	ValueList(ctx context.Context, req *api.ValueListRequest) (*api.ValueListResponse, error)
	// SetValue 设置脚本值
	SetValue(ctx context.Context, req *api.SetValueRequest) (*api.SetValueResponse, error)
	// DeleteValue 删除脚本值
	DeleteValue(ctx context.Context, req *api.DeleteValueRequest) (*api.DeleteValueResponse, error)
}

type valueSvc struct {
}

var defaultValue = &valueSvc{}

func Value() ValueSvc {
	return defaultValue
}

// ValueList 脚本值列表
func (v *valueSvc) ValueList(ctx context.Context, req *api.ValueListRequest) (*api.ValueListResponse, error) {
	scripts, err := script_repo.Script().FindByStorage(ctx, req.StorageName)
	if err != nil {
		return nil, err
	}
	if len(scripts) == 0 {
		return nil, i18n.NewNotFoundError(ctx, code.StorageNameNotFound)
	}
	script := scripts[0]
	list, _, err := value_repo.Value().FindPage(ctx, script.StorageName())
	if err != nil {
		return nil, err
	}
	resp := &api.ValueListResponse{
		List: make([]*api.Value, 0),
	}
	for _, v := range list {
		resp.List = append(resp.List, &api.Value{
			StorageName: v.StorageName,
			Key:         v.Key,
			Value:       v.Value,
			Createtime:  v.Createtime,
		})
	}
	return resp, nil
}

// SetValue 设置脚本值
func (v *valueSvc) SetValue(ctx context.Context, req *api.SetValueRequest) (*api.SetValueResponse, error) {
	scripts, err := script_repo.Script().FindByStorage(ctx, req.StorageName)
	if err != nil {
		return nil, err
	}
	if len(scripts) == 0 {
		return nil, i18n.NewNotFoundError(ctx, code.StorageNameNotFound)
	}
	script := scripts[0]
	for _, v := range req.Values {
		model, err := value_repo.Value().Find(ctx, script.StorageName(), v.Key)
		if err != nil {
			return nil, err
		}
		if model == nil {
			if err := value_repo.Value().Create(ctx, &value_entity.Value{
				StorageName: script.StorageName(),
				Key:         v.Key,
				Value:       v.Value,
				Createtime:  time.Now().Unix(),
			}); err != nil {
				return nil, err
			}
		} else {
			model.Value = v.Value
			if err := value_repo.Value().Update(ctx, model); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

// DeleteValue 删除脚本值
func (v *valueSvc) DeleteValue(ctx context.Context, req *api.DeleteValueRequest) (*api.DeleteValueResponse, error) {
	scripts, err := script_repo.Script().FindByStorage(ctx, req.StorageName)
	if err != nil {
		return nil, err
	}
	if len(scripts) == 0 {
		return nil, i18n.NewNotFoundError(ctx, code.StorageNameNotFound)
	}
	script := scripts[0]
	if err := value_repo.Value().Delete(ctx, script.StorageName(), req.Key); err != nil {
		return nil, err
	}
	return nil, nil
}
