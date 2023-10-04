package cloudcat_api

import (
	"context"

	"github.com/scriptscat/cloudcat/internal/api/scripts"
)

type Value struct {
	cli *Client
}

func NewValue(cli *Client) *Value {
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

func (s *Value) SetValue(ctx context.Context, req *scripts.SetValueRequest) (*scripts.SetValueResponse, error) {
	resp := &scripts.SetValueResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Value) DeleteValue(ctx context.Context, req *scripts.DeleteValueRequest) (*scripts.DeleteValueResponse, error) {
	resp := &scripts.DeleteValueResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}
