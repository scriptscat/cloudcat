package service

import (
	"github.com/scriptscat/cloudcat/internal/domain/sync/dto"
	"github.com/scriptscat/cloudcat/internal/domain/sync/entity"
	"github.com/scriptscat/cloudcat/internal/domain/sync/repository"
	"github.com/scriptscat/cloudcat/internal/pkg/cnt"
	"github.com/sirupsen/logrus"
)

type Sync interface {
	SyncScript(user, device int64, scripts []*dto.SyncScript) ([]*dto.SyncScript, error)
	FullScript(user, device int64) ([]*entity.SyncScript, error)

	UploadSubscribe(user, device int64)
	FullSubscribe(user, device int64) error

	SyncSetting(user, device int64) error
}

type sync struct {
	script repository.Script
}

func NewSync() Sync {
	return &sync{}
}

func (s *sync) SyncScript(user, device int64, scripts []*dto.SyncScript) ([]*dto.SyncScript, error) {
	ret := make([]*dto.SyncScript, len(scripts))
	for i, v := range scripts {
		script, err := s.script.FindByUUID(user, device, v.UUID)
		if err != nil {
			logrus.Warnf("sync script find script error: %v", err)
			ret[i] = &dto.SyncScript{
				Action: "error",
				Msg:    "同步失败,系统错误",
			}
			continue
		}
		if script != nil && script.Updatetime > v.Actiontime {
			// 更新时间大于操作时间,表明已经被操作过了,当前操作是过期操作
			if script.State == cnt.ACTIVE {
				ret[i] = &dto.SyncScript{
					Action:     "reinstall",
					Actiontime: script.Updatetime,
					UUID:       script.UUID,
					Script:     script,
				}
			} else {
				ret[i] = &dto.SyncScript{
					Action:     "uninstall",
					Actiontime: script.Updatetime,
					UUID:       script.UUID,
				}
			}
			continue
		}
		// 时间小或者是空,更新脚本
		v.Script.UserID = user
		v.Script.DeviceID = device
		v.Script.Updatetime = v.Actiontime
		if script != nil {
			v.Script.ID = script.ID
			v.Script.Createtime = script.Createtime
		}
		if err := s.script.Save(v.Script); err != nil {
			logrus.Warnf("sync script save script error: %v", err)
			ret[i] = &dto.SyncScript{
				Action: "error",
				Msg:    "同步失败,系统错误",
			}
			continue
		}
		ret[i] = &dto.SyncScript{
			Action: "ok",
		}
	}
	return ret, nil
}

func (s *sync) FullScript(user, device int64) ([]*entity.SyncScript, error) {
	list, err := s.script.ListScript(user, device)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *sync) UploadSubscribe(user, device int64) {

}

func (s *sync) FullSubscribe(user, device int64) error {
	panic("implement me")
}

func (s *sync) SyncSetting(user, device int64) error {
	panic("implement me")
}
