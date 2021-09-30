package entity

type SyncDevice struct {
	ID          int64  `gorm:"primaryKey" json:"-"`
	UserID      int64  `gorm:"index:device_user_id;column:user_id;type:bigint(20);not null" json:"user_id"`
	Name        string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Remark      string `gorm:"column:remark;type:varchar(255)" json:"remark"`
	Setting     string `gorm:"column:setting;type:mediumtext;not null" json:"setting"`
	Settingtime int64  `gorm:"column:settingtime;type:bigint(20);not null" json:"settingtime"`
	Createtime  int64  `gorm:"column:createtime;type:bigint(20);not null" json:"createtime"`
}
