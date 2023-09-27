package auth_ctr

import (
	"context"

	api "github.com/scriptscat/cloudcat/internal/api/auth"
	"github.com/scriptscat/cloudcat/internal/service/auth_svc"
)

type Token struct {
}

func NewToken() *Token {
	return &Token{}
}

// TokenList 获取token列表
func (t *Token) TokenList(ctx context.Context, req *api.TokenListRequest) (*api.TokenListResponse, error) {
	return auth_svc.Token().TokenList(ctx, req)
}

// TokenCreate 创建token
func (t *Token) TokenCreate(ctx context.Context, req *api.TokenCreateRequest) (*api.TokenCreateResponse, error) {
	return auth_svc.Token().TokenCreate(ctx, req)
}

// TokenDelete 删除token
func (t *Token) TokenDelete(ctx context.Context, req *api.TokenDeleteRequest) (*api.TokenDeleteResponse, error) {
	return auth_svc.Token().TokenDelete(ctx, req)
}
