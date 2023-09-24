package cloudcat_api

import (
	"context"

	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/api/scripts"
)

type Script struct {
	cli *mux.Client
}

func NewScript(cli *mux.Client) *Script {
	return &Script{
		cli: cli,
	}
}

func (s *Script) List(ctx context.Context, req *scripts.ListRequest) (*scripts.ListResponse, error) {
	resp := &scripts.ListResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Script) Install(ctx context.Context, req *scripts.InstallRequest) (*scripts.InstallResponse, error) {
	resp := &scripts.InstallResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Script) Get(ctx context.Context, req *scripts.GetRequest) (*scripts.GetResponse, error) {
	resp := &scripts.GetResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Script) Update(ctx context.Context, req *scripts.UpdateRequest) (*scripts.UpdateResponse, error) {
	resp := &scripts.UpdateResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Script) Delete(ctx context.Context, req *scripts.DeleteRequest) (*scripts.DeleteResponse, error) {
	resp := &scripts.DeleteResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}
