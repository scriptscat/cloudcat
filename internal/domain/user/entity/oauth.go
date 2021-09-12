package entity

type BbsOauthUser struct {
	ID         int64  `gorm:"primaryKey;column:id;type:bigint(20);not null" json:"-"`
	Openid     string `gorm:"unique;column:openid;type:varchar(255);not null" json:"openid"`
	UserID     int64  `gorm:"index:bbs_user_id;column:user_id;type:bigint(20);not null" json:"user_id"`
	Status     int8   `gorm:"column:status;type:tinyint(4);not null;default:1" json:"status"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20);not null" json:"createtime"`
}

type WechatOauthUser struct {
	ID         int64  `gorm:"primaryKey;column:id;type:bigint(20);not null" json:"-"`
	Openid     string `gorm:"unique;column:openid;type:varchar(255);not null" json:"openid"`
	Unionid    string `gorm:"unique;column:unionid;type:varchar(255)" json:"unionid"`
	UserID     int64  `gorm:"index:wc_user_id;column:user_id;type:bigint(20);not null" json:"user_id"`
	Status     int8   `gorm:"column:status;type:tinyint(4);not null" json:"status"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20);not null" json:"createtime"`
}
