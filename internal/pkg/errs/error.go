package errs

type JsonRespondError struct {
	Status int    `json:"-"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

func NewError(status, code int, msg string) error {
	return &JsonRespondError{
		Status: status,
		Code:   code,
		Msg:    msg,
	}
}

func (j *JsonRespondError) Error() string {
	return j.Msg
}
