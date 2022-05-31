// Author: chenqionghe
// Time: 2018-10
// 提供加解密功能

package base64

import (
	"encoding/base64"
)

// base64Encode
func Base64Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// base64Decode
func Base64Decode(value string) string {
	result, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return ""
	}
	return string(result)
}
