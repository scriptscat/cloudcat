package entity

type SyncSetting struct {
	ID         int64  `gorm:"primaryKey" json:"-"`
	UserID     int64  `gorm:"uniqueIndex:user_device;index:sync_setting_user_id;column:user_id;type:bigint(20);not null" json:"user_id"`
	DeviceID   int64  `gorm:"uniqueIndex:user_device;column:device_id;type:bigint(20);not null" json:"device_id"`
	Setting    string `gorm:"column:setting;type:mediumtext;not null" json:"setting"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20);not null" json:"createtime"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20);not null" json:"updatetime"`
}

type SyncValue struct {
	UserID      int64
	DeviceID    int64
	ScriptUUID  string
	StorageName string
	Key         string
	Value       string
}
