package auth_ctr

import (
	"context"

	api "github.com/scriptscat/cloudcat/internal/api/auth"
	"github.com/scriptscat/cloudcat/internal/service/auth_svc"
)

type Secret struct {
}

func NewSecret() *Secret {
	return &Secret{}
}

// TokenList 获取token列表
func (s *Secret) TokenList(ctx context.Context, req *api.TokenListRequest) (*api.TokenListResponse, error) {
	return auth_svc.Secret().TokenList(ctx, req)
}

// TokenCreate 创建token
func (s *Secret) TokenCreate(ctx context.Context, req *api.TokenCreateRequest) (*api.TokenCreateResponse, error) {
	return auth_svc.Secret().TokenCreate(ctx, req)
}

// TokenDelete 删除token
func (s *Secret) TokenDelete(ctx context.Context, req *api.TokenDeleteRequest) (*api.TokenDeleteResponse, error) {
	return auth_svc.Secret().TokenDelete(ctx, req)
}
