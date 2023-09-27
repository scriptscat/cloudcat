package cloudcat_api

import (
	"context"
	"github.com/scriptscat/cloudcat/internal/api/auth"
)

type Token struct {
	cli *Client
}

func NewToken(cli *Client) *Token {
	return &Token{
		cli: cli,
	}
}

func (t *Token) Create(ctx context.Context, req *auth.TokenCreateRequest) (*auth.TokenCreateResponse, error) {
	resp := &auth.TokenCreateResponse{}
	if err := t.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (t *Token) List(ctx context.Context, req *auth.TokenListRequest) (*auth.TokenListResponse, error) {
	resp := &auth.TokenListResponse{}
	if err := t.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (t *Token) Delete(ctx context.Context, req *auth.TokenDeleteRequest) (*auth.TokenDeleteResponse, error) {
	resp := &auth.TokenDeleteResponse{}
	if err := t.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}
