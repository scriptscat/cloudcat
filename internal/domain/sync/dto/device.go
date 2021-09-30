package dto

import "github.com/scriptscat/cloudcat/internal/domain/sync/entity"

type Device struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Name        string `json:"name"`
	Remark      string `json:"remark"`
	Setting     string `json:"setting"`
	Settingtime int64  `json:"settingtime"`
	Createtime  int64  `json:"createtime"`
	SyncVersion struct {
		Script    int64 `json:"script"`
		Subscribe int64 `json:"subscribe"`
	} `json:"sync_version"`
}

func ToDevice(device *entity.SyncDevice, script, subscribe int64) *Device {
	return &Device{
		ID:          device.ID,
		UserID:      device.UserID,
		Name:        device.Name,
		Remark:      device.Remark,
		Setting:     device.Setting,
		Settingtime: device.Settingtime,
		Createtime:  device.Createtime,
		SyncVersion: struct {
			Script    int64 `json:"script"`
			Subscribe int64 `json:"subscribe"`
		}{
			Script:    script,
			Subscribe: subscribe,
		},
	}
}
