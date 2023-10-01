package value_entity

import (
	"encoding/json"
)

type Value struct {
	StorageName string      `json:"storage_name"`
	Key         string      `json:"key"`
	Value       ValueString `json:"value"`
	Createtime  int64       `json:"createtime"`
}

type ValueString struct {
	value []byte
}

func (v *ValueString) Set(value interface{}) error {
	var err error
	v.value, err = json.Marshal(value)
	if err != nil {
		return err
	}
	return nil
}

func (v *ValueString) Get() interface{} {
	var value interface{}
	err := json.Unmarshal(v.value, &value)
	if err != nil {
		return nil
	}
	return value
}

func (v *ValueString) MarshalJSON() ([]byte, error) {
	return v.value, nil
}

func (v *ValueString) UnmarshalJSON(data []byte) error {
	v.value = data
	return nil
}
