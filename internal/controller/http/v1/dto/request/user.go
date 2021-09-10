package request

type UserLogin struct {
	Username string `form:"name" binding:"omitempty,alphanumunicode,min=3,max=6" label:"库的名字"`
	Email    string `form:"email" binding:"omitempty,email" label:"用户邮箱"`
	Password string `form:"password" binding:"required,min=6,max=18" label:"用户密码"`
}
