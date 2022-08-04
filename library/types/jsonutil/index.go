package jsonutil

import (
	jsoniter "github.com/json-iterator/go"
)

func ToBytes(i interface{}) []byte {
	bytes, err := jsoniter.Marshal(i)
	if err != nil {
		return nil
	}
	return bytes
}

func ToString(i interface{}) string {
	switch vv := i.(type) {
	case []byte:
		return string(vv)
	default:
		bytes, err := jsoniter.Marshal(i)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
}

func ToJson(i, v interface{}) (err error) {
	switch vv := i.(type) {
	case string:
		if err = jsoniter.Unmarshal([]byte(vv), v); err != nil {
			return
		}
	case []byte:
		if err = jsoniter.Unmarshal(vv, v); err != nil {
			return
		}
	default:
		if err = jsoniter.Unmarshal(ToBytes(i), v); err != nil {
			return
		}
	}
	return
}
