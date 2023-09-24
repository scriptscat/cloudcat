package auth_svc

import (
	"context"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/scriptscat/cloudcat/internal/model/entity/token_entity"
	"github.com/scriptscat/cloudcat/internal/repository/token_repo"
	"time"

	api "github.com/scriptscat/cloudcat/internal/api/auth"
)

type SecretSvc interface {
	// TokenList 获取token列表
	TokenList(ctx context.Context, req *api.TokenListRequest) (*api.TokenListResponse, error)
	// TokenCreate 创建token
	TokenCreate(ctx context.Context, req *api.TokenCreateRequest) (*api.TokenCreateResponse, error)
	// TokenDelete 删除token
	TokenDelete(ctx context.Context, req *api.TokenDeleteRequest) (*api.TokenDeleteResponse, error)
}

type secretSvc struct {
}

var defaultSecret = &secretSvc{}

func Secret() SecretSvc {
	return defaultSecret
}

// TokenList 获取token列表
func (s *secretSvc) TokenList(ctx context.Context, req *api.TokenListRequest) (*api.TokenListResponse, error) {
	list, err := token_repo.Secret().FindPage(ctx)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// TokenCreate 创建token
func (s *secretSvc) TokenCreate(ctx context.Context, req *api.TokenCreateRequest) (*api.TokenCreateResponse, error) {
	id := utils.RandString(16, utils.Letter)
	secret := utils.RandString(32, utils.Letter)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	m := &token_entity.Token{
		ID:                id,
		Token:             tokenString,
		Secret:            secret,
		DataEncryptionKey: utils.RandString(32, utils.Letter),
		Status:            consts.ACTIVE,
		Createtime:        time.Now().Unix(),
		Updatetime:        time.Now().Unix(),
	}
	if err := token_repo.Secret().Create(ctx, m); err != nil {
		return nil, err
	}
	return &api.TokenCreateResponse{
		Token: &api.Token{
			ID:                m.ID,
			Token:             m.Token,
			DataEncryptionKey: m.DataEncryptionKey,
			Status:            m.Status,
			Createtime:        m.Createtime,
			Updatetime:        m.Createtime,
		}}, nil
}

// TokenDelete 删除token
func (s *secretSvc) TokenDelete(ctx context.Context, req *api.TokenDeleteRequest) (*api.TokenDeleteResponse, error) {
	err := token_repo.Secret().Delete(ctx, req.TokenID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
