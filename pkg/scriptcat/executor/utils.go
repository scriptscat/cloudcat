package executor

import (
	"rogchap.com/v8go"
)

func getString(v *v8go.Value) string {
	if !v.IsString() {
		return ""
	}
	return v.String()
}

func getObjString(obj *v8go.Object, key string) string {
	ret, err := obj.Get(key)
	if err != nil {
		return ""
	}
	if !ret.IsString() {
		return ""
	}
	return ret.String()
}

func getObjBool(obj *v8go.Object, key string) bool {
	ret, err := obj.Get(key)
	if err != nil {
		return false
	}
	if !ret.IsBoolean() {
		return false
	}
	return ret.Boolean()
}

func getFunction(obj *v8go.Object, key string) *v8go.Function {
	ret, err := obj.Get(key)
	if err != nil {
		return nil
	}
	f, err := ret.AsFunction()
	if err != nil {
		return nil
	}
	return f
}

func getObject(obj *v8go.Object, key string) *v8go.Object {
	ret, err := obj.Get(key)
	if err != nil {
		return nil
	}
	o, err := ret.AsObject()
	if err != nil {
		return nil
	}
	return o
}

func getNumber(obj *v8go.Object, key string) float64 {
	ret, err := obj.Get(key)
	if err != nil {
		return 0
	}
	if !ret.IsNumber() {
		return 0
	}
	return ret.Number()
}
