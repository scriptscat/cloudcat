package entity

type User struct {
	ID           int64  `gorm:"primaryKey;column:id;type:bigint(20);not null" json:"-"`                    // 用户id
	Username     string `gorm:"index:username;column:username;type:varchar(255);not null" json:"username"` // 用户名
	PasswordHash string `gorm:"column:password_hash;type:varchar(255)" json:"password_hash"`
	Email        string `gorm:"index:email;column:email;type:varchar(255)" json:"email"`
	Mobile       string `gorm:"index:mobile;column:mobile;type:varchar(255)" json:"mobile"`
	Createtime   int64  `gorm:"column:createtime;type:bigint(20);not null" json:"createtime"`
	Updatetime   int64  `gorm:"column:updatetime;type:bigint(20)" json:"updatetime"`
}
