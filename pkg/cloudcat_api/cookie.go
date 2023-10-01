package cloudcat_api

import (
	"context"

	"github.com/scriptscat/cloudcat/internal/api/scripts"
)

type Cookie struct {
	cli *Client
}

func NewCookie(cli *Client) *Cookie {
	return &Cookie{
		cli: cli,
	}
}

func (s *Cookie) CookieList(ctx context.Context, req *scripts.CookieListRequest) (*scripts.CookieListResponse, error) {
	resp := &scripts.CookieListResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Cookie) DeleteCookie(ctx context.Context, req *scripts.DeleteCookieRequest) (*scripts.DeleteCookieResponse, error) {
	resp := &scripts.DeleteCookieResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Cookie) SetCookie(ctx context.Context, req *scripts.SetCookieRequest) (*scripts.SetCookieResponse, error) {
	resp := &scripts.SetCookieResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}
