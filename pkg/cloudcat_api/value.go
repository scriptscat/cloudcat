package cloudcat_api

import (
	"context"

	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/api/scripts"
)

type Value struct {
	cli *mux.Client
}

func NewValue(cli *mux.Client) *Value {
	return &Value{
		cli: cli,
	}
}

func (s *Value) ValueList(ctx context.Context, req *scripts.ValueListRequest) (*scripts.ValueListResponse, error) {
	resp := &scripts.ValueListResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}
