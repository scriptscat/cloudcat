package auth

import (
	"github.com/codfrm/cago/server/mux"
)

type Token struct {
	ID                string `json:"id"`
	Token             string `json:"token"`
	DataEncryptionKey string `json:"data_encryption_key" yaml:"dataEncryptionKey"`
	Createtime        int64  `json:"createtime"`
	Updatetime        int64  `json:"updatetime"`
}

type TokenListRequest struct {
	mux.Meta `path:"/tokens" method:"GET"`
	TokenID  string `form:"token_id"`
}

type TokenListResponse struct {
	List []*Token `json:"list"`
}

type TokenCreateRequest struct {
	mux.Meta `path:"/tokens" method:"POST"`
	TokenID  string `form:"token_id" json:"token_id"`
}

type TokenCreateResponse struct {
	Token *Token `json:"token"`
}

type TokenDeleteRequest struct {
	mux.Meta `path:"/tokens/:tokenId" method:"DELETE"`
	TokenID  string `uri:"tokenId"`
}

type TokenDeleteResponse struct {
}
