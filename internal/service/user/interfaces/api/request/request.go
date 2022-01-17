package request

type ResetPasswordRequest struct {
	Code       string `form:"code" binding:"required" label:"验证码"`
	Password   string `form:"password" binding:"required,min=6,max=18" label:"密码"`
	Repassword string `form:"repassword" binding:"required,min=6,max=18,eqfield=Password" label:"再输入一次密码"`
}
