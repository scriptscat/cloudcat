package entity

type SyncScript struct {
	ID           int64  `gorm:"primaryKey" json:"-"`
	UserID       int64  `gorm:"index:sync_script_user_id;column:user_id;type:bigint(20);not null" json:"user_id"`
	DeviceID     int64  `gorm:"uniqueIndex:device_uuid;column:device_id;type:bigint(20);not null" json:"device_id"`
	Name         string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	UUID         string `gorm:"uniqueIndex:device_uuid;column:uuid;type:varchar(128);not null" json:"uuid"`
	Code         string `gorm:"column:code;type:mediumtext;not null" json:"code"`
	MetaJSON     string `gorm:"column:meta_json;type:text;not null" json:"meta_json"`
	SelfMeta     string `gorm:"column:self_meta;type:text" json:"self_meta"`
	Origin       string `gorm:"column:origin;type:text;not null" json:"origin"`
	Sort         int32  `gorm:"column:sort;type:int(10);default:0" json:"sort"`
	SubscribeUrl string `gorm:"column:subscribe_url;type:text;" json:"subscribe_url"`
	Type         int8   `gorm:"column:type;type:tinyint(4);not null" json:"type"`
	Status       int8   `gorm:"column:status;type:tinyint(4);not null" json:"status"`
	Createtime   int64  `gorm:"column:createtime;type:bigint(20);not null" json:"createtime"`
	Updatetime   int64  `gorm:"column:updatetime;type:bigint(20);not null" json:"updatetime"`
}

type SyncSubscribe struct {
	ID         int64  `gorm:"primaryKey" json:"-"`
	UserID     int64  `gorm:"index:sync_subscribe_user_id;column:user_id;type:bigint(20);not null" json:"user_id"`
	DeviceID   int64  `gorm:"uniqueIndex:device_url;column:device_id;type:bigint(20);not null" json:"device_id"`
	Name       string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	URL        string `gorm:"column:url;type:text;not null" json:"url"`
	Code       string `gorm:"column:code;type:text;not null" json:"code"`
	MetaJSON   string `gorm:"column:meta_json;type:text;not null" json:"meta_json"`
	Scripts    string `gorm:"column:scripts;type:text;not null" json:"scripts"`
	Status     int8   `gorm:"column:status;type:tinyint(4);not null" json:"status"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20);not null" json:"createtime"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20);not null" json:"updatetime"`
}
