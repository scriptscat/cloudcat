package value_entity

import "time"

type Value struct {
	StorageName string    `json:"storage_name"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	CreatedTime time.Time `json:"created_time"`
}
