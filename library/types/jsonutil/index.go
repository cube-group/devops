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
	bytes, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func ToJson(i string) gin.H {
	if i == "" {
		return gin.H{}
	}
	var res gin.H
	err := json.Unmarshal([]byte(i), &res)
	if err != nil {
		return gin.H{}
	}
	return res
}
