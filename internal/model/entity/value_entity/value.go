package value_entity

type Value struct {
	StorageName string `json:"storage_name"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Createtime  int64  `json:"createtime"`
}
