package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/scriptscat/cloudcat/internal/domain/sync/dto"
	"github.com/scriptscat/cloudcat/internal/domain/sync/entity"
	"github.com/scriptscat/cloudcat/internal/domain/sync/repository"
	mock_repository "github.com/scriptscat/cloudcat/internal/domain/sync/repository/mock"
)

func Test_sync_PushScript(t *testing.T) {
	type args struct {
		user    int64
		device  int64
		version int64
		scripts []*dto.SyncScript
	}
	success := func(mockctl *gomock.Controller, args args) repository.Script {
		mock := mock_repository.NewMockScript(mockctl)
		mock.EXPECT().LatestVersion(args.user, args.device).Return(args.version, nil).Times(1)

		l := 0
		for _, v := range args.scripts {
			if v.Msg == "nil" {
				l++
				mock.EXPECT().FindByUUID(args.user, args.device, v.UUID).Return(nil, nil)
			} else if v.Msg != "" {
				mock.EXPECT().FindByUUID(args.user, args.device, v.UUID).Return(nil, errors.New(v.Msg))
			} else {
				l++
				mock.EXPECT().FindByUUID(args.user, args.device, v.UUID).Return(v.Script, nil)
			}
		}

		mock.EXPECT().Save(gomock.Any()).Return(nil).AnyTimes()

		mock.EXPECT().PushVersion(args.user, args.device, gomock.Len(l)).Return(args.version+1, nil).Times(1)

		return mock
	}
	nils := &entity.SyncScript{}
	id1 := &entity.SyncScript{ID: 1}
	id2 := &entity.SyncScript{ID: 2}
	tests := []struct {
		name    string
		fields  func(*gomock.Controller, args) repository.Script
		args    args
		want    []*dto.SyncScript
		want1   int64
		wantErr bool
	}{
		{"版本号不等错误", func(mockctl *gomock.Controller, args args) repository.Script {
			mock := mock_repository.NewMockScript(mockctl)
			mock.EXPECT().LatestVersion(args.user, args.device).Return(args.version+1, nil)
			return mock
		}, args{
			user:    1,
			device:  1,
			version: 1,
			scripts: []*dto.SyncScript{{}},
		}, nil, 0, true},
		{"空的推送", func(controller *gomock.Controller, args args) repository.Script {
			return nil
		}, args{
			user:    1,
			device:  1,
			version: 1,
			scripts: []*dto.SyncScript{},
		}, nil, 0, true},
		{"单条", success, args{
			user:    2,
			device:  2,
			version: 2,
			scripts: []*dto.SyncScript{{
				Action:     "reinstall",
				Actiontime: 1,
				UUID:       "uuid1",
				Script:     nils,
			}},
		}, []*dto.SyncScript{{Action: "ok", UUID: "uuid1", Script: nils}}, 3, false},
		{"多条覆盖全部条件", success, args{
			user:    2,
			device:  2,
			version: 2,
			scripts: []*dto.SyncScript{{
				Action:     "reinstall",
				Actiontime: 1,
				UUID:       "uuid1",
				Script:     nils,
			}, {
				Action:     "uninstall",
				Actiontime: 1,
				UUID:       "uuid2",
				Msg:        "error",
				Script:     nils,
			}, {
				Action:     "install",
				Actiontime: 1,
				UUID:       "uuid3",
				Script:     id1,
			}, {
				Action:     "uninstall",
				Actiontime: 1,
				UUID:       "uuid4",
				Msg:        "nil",
				Script:     id2,
			}},
		}, []*dto.SyncScript{
			{Action: "ok", UUID: "uuid1", Script: nils},
			{Action: "error", UUID: "", Msg: "同步失败,系统错误"},
			{Action: "ok", UUID: "uuid3", Script: id1},
			{Action: "ok", UUID: "uuid4", Script: id2},
		}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockctl := gomock.NewController(t)
			defer mockctl.Finish()
			mock := mock_repository.NewMockDevice(mockctl)
			mock.EXPECT().FindById(gomock.Any()).Return(&entity.SyncDevice{
				ID:     tt.args.device,
				UserID: tt.args.user,
			}, nil).AnyTimes()
			subMock := mock_repository.NewMockSubscribe(mockctl)
			subMock.EXPECT().LatestVersion(gomock.Any(), gomock.Any()).Return(int64(1), nil).AnyTimes()
			s := &sync{
				device:    mock,
				subscribe: subMock,
				script:    tt.fields(mockctl, tt.args),
			}
			got, got1, err := s.PushScript(tt.args.user, tt.args.device, tt.args.version, tt.args.scripts)
			if (err != nil) != tt.wantErr {
				t.Errorf("PushScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PushScript() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PushScript() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_sync_PullScript(t *testing.T) {
	type fields struct {
		script repository.Script
	}
	type args struct {
		user    int64
		device  int64
		version int64
	}
	tests := []struct {
		name    string
		fields  func(mockctl *gomock.Controller, args args) repository.Script
		args    args
		want    []*dto.SyncScript
		want1   int64
		wantErr bool
	}{
		{"版本相等", func(mockctl *gomock.Controller, args args) repository.Script {
			s := mock_repository.NewMockScript(mockctl)
			s.EXPECT().LatestVersion(args.user, args.device).Return(args.version, nil)
			return s
		}, args{
			user:    1,
			device:  1,
			version: 1,
		}, nil, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockctl := gomock.NewController(t)
			defer mockctl.Finish()
			s := &sync{
				script: tt.fields(mockctl, tt.args),
			}
			got, got1, err := s.PullScript(tt.args.user, tt.args.device, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("PullScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullScript() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PullScript() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
