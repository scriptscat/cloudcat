package bbs

type ErrorRespond struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *ErrorRespond) Error() string {
	return e.Msg
}

type AccessTokenRespond struct {
	ErrorRespond
	AccessToken string `json:"access_token"`
}

type UserRespond struct {
	ErrorRespond
	User struct {
		Uid      string `json:"uid"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Email    string `json:"email"`
	} `json:"user"`
}
