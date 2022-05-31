// Author: chenqionghe
// Time: 2018-10
// 提供加解密功能

package sha1

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// 生成小写Sha1
func Sha1(text string) string {
	ctx := sha1.New()
	_, err := ctx.Write([]byte(text))
	if err != nil {
		return ""
	}
	return strings.ToLower(hex.EncodeToString(ctx.Sum(nil)))
}

// 生成大写Sha1
func Sha1Upper(text string) string {
	result := Sha1(text)
	return strings.ToUpper(result)
}
