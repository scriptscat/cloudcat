package cloudcat_api

import (
	"context"

	"github.com/scriptscat/cloudcat/internal/api/scripts"
)

type Script struct {
	cli *Client
}

func NewScript(cli *Client) *Script {
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

func (s *Script) Run(ctx context.Context, req *scripts.RunRequest) (*scripts.RunResponse, error) {
	resp := &scripts.RunResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Script) Stop(ctx context.Context, req *scripts.StopRequest) (*scripts.StopResponse, error) {
	resp := &scripts.StopResponse{}
	if err := s.cli.Do(ctx, req, resp); err != nil {
		return resp, err
	}
	return resp, nil
}
