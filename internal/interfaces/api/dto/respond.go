package dto

type List struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
