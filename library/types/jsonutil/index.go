package jsonutil

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func ToBytes(i interface{}) []byte {
	bytes, err := json.Marshal(i)
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
		bytes, err := json.Marshal(i)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
}

func ToJson(i interface{}) (res gin.H) {
	switch vv := i.(type) {
	case string:
		err := json.Unmarshal([]byte(vv), &res)
		if err != nil {
			return gin.H{}
		}
	case []byte:
		err := json.Unmarshal(vv, &res)
		if err != nil {
			return gin.H{}
		}
	}
	return
}
