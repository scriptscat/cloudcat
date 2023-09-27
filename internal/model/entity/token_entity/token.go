package token_entity

type Token struct {
	ID                string `json:"id"`
	Token             string `json:"token"`               // jwt
	Secret            string `json:"secret"`              // jwt签名密钥
	DataEncryptionKey string `json:"data_encryption_key"` // 数据加密密钥
	Status            int8   `json:"status"`
	Createtime        int64  `json:"createtime"`
	Updatetime        int64  `json:"updatetime"`
}
