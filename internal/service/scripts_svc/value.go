package scripts_svc

import (
	"context"
	"github.com/scriptscat/cloudcat/internal/repository/script_repo"

	api "github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/internal/repository/value_repo"
)

type ValueSvc interface {
	// ValueList 脚本值列表
	ValueList(ctx context.Context, req *api.ValueListRequest) (*api.ValueListResponse, error)
	// StorageList 值储存空间列表
	StorageList(ctx context.Context, req *api.StorageListRequest) (*api.StorageListResponse, error)
}

type valueSvc struct {
}

var defaultValue = &valueSvc{}

func Value() ValueSvc {
	return defaultValue
}

// ValueList 脚本值列表
func (v *valueSvc) ValueList(ctx context.Context, req *api.ValueListRequest) (*api.ValueListResponse, error) {
	list, _, err := value_repo.Value().FindPage(ctx, req.StorageName)
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
			CreatedTime: v.CreatedTime,
		})
	}
	return resp, nil
}

// StorageList 值储存空间列表
func (v *valueSvc) StorageList(ctx context.Context, req *api.StorageListRequest) (*api.StorageListResponse, error) {
	list, err := script_repo.Script().StorageList(ctx)
	if err != nil {
		return nil, err
	}
	resp := &api.StorageListResponse{
		List: make([]*api.Storage, 0),
	}
	for _, v := range list {
		resp.List = append(resp.List, &api.Storage{
			Name:         v.Name,
			LinkScriptID: v.LinkScriptID,
		})
	}
	return resp, nil
}
