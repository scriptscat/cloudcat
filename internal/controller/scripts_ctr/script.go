package scripts_ctr

import (
	"context"

	api "github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/internal/service/scripts_svc"
)

type Script struct {
}

func NewScripts() *Script {
	return &Script{}
}

// List 脚本列表
func (s *Script) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return scripts_svc.Script().List(ctx, req)
}

// Install 安装脚本
func (s *Script) Install(ctx context.Context, req *api.InstallRequest) (*api.InstallResponse, error) {
	return scripts_svc.Script().Install(ctx, req)
}

// Get 获取脚本
func (s *Script) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	return scripts_svc.Script().Get(ctx, req)
}

// Update 更新脚本
func (s *Script) Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponse, error) {
	return scripts_svc.Script().Update(ctx, req)
}

// Delete 删除脚本
func (s *Script) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	return scripts_svc.Script().Delete(ctx, req)
}

// StorageList 值储存空间列表
func (s *Script) StorageList(ctx context.Context, req *api.StorageListRequest) (*api.StorageListResponse, error) {
	return scripts_svc.Script().StorageList(ctx, req)
}
