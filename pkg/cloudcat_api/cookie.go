package cloudcat_api

import (
	"context"

	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/api/scripts"
)

type Cookie struct {
	cli *mux.Client
}

func NewCookie(cli *mux.Client) *Cookie {
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
