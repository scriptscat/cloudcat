package application

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/entity"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/repository"
	mock_repository "github.com/scriptscat/cloudcat/internal/service/user/domain/repository/mock"
	"github.com/scriptscat/cloudcat/internal/service/user/domain/vo"
)

func Test_oauth_OAuthPlatform(t *testing.T) {
	type args struct {
		uid int64
	}
	tests := []struct {
		name    string
		fields  func(mockctl *gomock.Controller, args args) (repository.BBSOAuth, repository.WechatOAuth)
		args    args
		want    *vo.OpenPlatform
		wantErr bool
	}{
		{name: "正常通过", fields: func(mockctl *gomock.Controller, args args) (repository.BBSOAuth, repository.WechatOAuth) {
			b, w := mock_repository.NewMockBBSOAuth(mockctl), mock_repository.NewMockWechatOAuth(mockctl)
			b.EXPECT().FindByUid(args.uid).Return(&entity.BbsOauthUser{}, nil)
			w.EXPECT().FindByUid(args.uid).Return(&entity.WechatOauthUser{}, nil)
			return b, w
		}, args: args{uid: 1}, want: &vo.OpenPlatform{Bbs: true, Wechat: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockctl := gomock.NewController(t)
			b, w := tt.fields(mockctl, tt.args)
			o := &oauth{
				bbsOAuthRepo:    b,
				wechatOAuthRepo: w,
			}
			got, err := o.OAuthPlatform(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("OAuthPlatform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OAuthPlatform() got = %v, want %v", got, tt.want)
			}
		})
	}
}
