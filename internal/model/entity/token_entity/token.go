package token_entity

type Token struct {
	ID                string `json:"id"`
	Token             string `json:"token"`
	Secret            string `json:"secret"`
	DataEncryptionKey string `json:"data_encryption_key"`
	Status            int8   `json:"status"`
	Createtime        int64  `json:"createtime"`
	Updatetime        int64  `json:"updatetime"`
}
