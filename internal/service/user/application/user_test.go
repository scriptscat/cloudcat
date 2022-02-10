package application

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/scriptscat/cloudcat/internal/infrastructure/config"
	mock_config "github.com/scriptscat/cloudcat/internal/infrastructure/config/mock"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/errs"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/repository"
	mock_repository "github.com/scriptscat/cloudcat/internal/service/user/domain/repository/mock"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/vo"
)

func Test_user_UpdateUserInfo(t *testing.T) {
	type args struct {
		uid int64
		req *vo.UpdateUserInfo
	}
	tests := []struct {
		name    string
		fields  func(mockctl *gomock.Controller, args args) repository.User
		args    args
		wantErr bool
	}{
		{"正常通过", func(mockctl *gomock.Controller, args args) repository.User {
			mock := mock_repository.NewMockUser(mockctl)
			mock.EXPECT().FindById(args.uid).Return(&entity.User{ID: 1}, nil)
			mock.EXPECT().FindByName(args.req.Username).Return(nil, errs.ErrUserNotFound)
			mock.EXPECT().SaveUser(gomock.Any()).Times(1)
			return mock
		}, args{uid: 1, req: &vo.UpdateUserInfo{Username: "用户名"}}, false},
		{
			"用户不存在", func(mockctl *gomock.Controller, args args) repository.User {
				mock := mock_repository.NewMockUser(mockctl)
				mock.EXPECT().FindById(args.uid).Return(nil, errs.ErrUserNotFound)
				return mock
			}, args{uid: 2, req: nil}, true,
		}, {
			"用户名存在", func(mockctl *gomock.Controller, args args) repository.User {
				mock := mock_repository.NewMockUser(mockctl)
				mock.EXPECT().FindById(args.uid).Return(&entity.User{ID: 1}, nil)
				mock.EXPECT().FindByName(args.req.Username).Return(&entity.User{ID: 3}, nil)
				return mock
			}, args{uid: 1, req: &vo.UpdateUserInfo{Username: "用户名"}}, true,
		}, {
			"系统错误", func(mockctl *gomock.Controller, args args) repository.User {
				mock := mock_repository.NewMockUser(mockctl)
				mock.EXPECT().FindById(args.uid).Return(&entity.User{ID: 1}, nil)
				mock.EXPECT().FindByName(args.req.Username).Return(nil, errors.New("system error"))
				return mock
			}, args{uid: 1, req: &vo.UpdateUserInfo{Username: "用户名"}}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockctl := gomock.NewController(t)
			u := &user{
				userRepo: tt.fields(mockctl, tt.args),
			}
			if err := u.UpdateUserInfo(tt.args.uid, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_user_UpdateEmail(t *testing.T) {
	type args struct {
		uid int64
		req *vo.UpdateEmail
	}
	tests := []struct {
		name    string
		fields  func(mockctl *gomock.Controller, args args) (config.SystemConfig, repository.User, repository.VerifyCode)
		args    args
		wantErr bool
	}{
		{"正常通过", func(mockctl *gomock.Controller, args args) (config.SystemConfig, repository.User, repository.VerifyCode) {
			s := mock_config.NewMockSystemConfig(mockctl)
			u := mock_repository.NewMockUser(mockctl)
			v := mock_repository.NewMockVerifyCode(mockctl)
			s.EXPECT().GetConfig(AllowEmailSuffix).Return("admin.com", nil)
			u.EXPECT().FindById(args.uid).Return(&entity.User{ID: 1}, nil)
			u.EXPECT().FindByEmail(args.req.Email).Return(nil, errs.ErrUserNotFound)
			v.EXPECT().FindById(args.req.Email).Return(&entity.VerifyCode{
				Identifier: args.req.Email,
				Op:         "change-user-email",
				Code:       "vcode1",
				Expired:    time.Now().Add(time.Second).Unix(),
			}, nil)
			v.EXPECT().InvalidCode(gomock.Any()).Return(nil)
			u.EXPECT().SaveUser(gomock.Any()).Times(1)
			return s, u, v
		}, args{uid: 1, req: &vo.UpdateEmail{
			Email: "admin@admin.com",
			Code:  "vcode1",
		}}, false},
		{name: "邮箱相等", fields: func(mockctl *gomock.Controller, args args) (config.SystemConfig, repository.User, repository.VerifyCode) {
			s := mock_config.NewMockSystemConfig(mockctl)
			u := mock_repository.NewMockUser(mockctl)
			v := mock_repository.NewMockVerifyCode(mockctl)
			u.EXPECT().FindById(args.uid).Return(&entity.User{ID: 1, Email: sql.NullString{String: args.req.Email, Valid: true}}, nil)
			return s, u, v
		}, args: args{uid: 1, req: &vo.UpdateEmail{
			Email: "admin@admin.com",
			Code:  "vcode1",
		}}, wantErr: true},
		{"邮箱存在", func(mockctl *gomock.Controller, args args) (config.SystemConfig, repository.User, repository.VerifyCode) {
			s := mock_config.NewMockSystemConfig(mockctl)
			u := mock_repository.NewMockUser(mockctl)
			v := mock_repository.NewMockVerifyCode(mockctl)
			s.EXPECT().GetConfig(AllowEmailSuffix).Return("admin.com", nil)
			u.EXPECT().FindById(args.uid).Return(&entity.User{ID: 1}, nil)
			u.EXPECT().FindByEmail(args.req.Email).Return(&entity.User{ID: 2}, nil)
			return s, u, v
		}, args{uid: 1, req: &vo.UpdateEmail{
			Email: "admin@admin.com",
			Code:  "vcode1",
		}}, true},
		{"后缀不允许", func(mockctl *gomock.Controller, args args) (config.SystemConfig, repository.User, repository.VerifyCode) {
			s := mock_config.NewMockSystemConfig(mockctl)
			u := mock_repository.NewMockUser(mockctl)
			v := mock_repository.NewMockVerifyCode(mockctl)
			s.EXPECT().GetConfig(AllowEmailSuffix).Return("admin.com,badmin.com", nil)
			u.EXPECT().FindById(args.uid).Return(&entity.User{ID: 1}, nil)
			return s, u, v
		}, args{uid: 1, req: &vo.UpdateEmail{
			Email: "admin@hadmin.com",
			Code:  "vcode1",
		}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockctl := gomock.NewController(t)
			c, us, v := tt.fields(mockctl, tt.args)
			u := &user{
				config:     c,
				userRepo:   us,
				verifyRepo: v,
			}
			if err := u.UpdateEmail(tt.args.uid, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
