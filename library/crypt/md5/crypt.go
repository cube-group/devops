// Author: chenqionghe
// Time: 2018-10
// 提供加解密功能

package md5

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// 32位小写MD5
func MD5(text string) string {
	ctx := md5.New()
	_, err := ctx.Write([]byte(text))
	if err != nil {
		return ""
	}
	return strings.ToLower(hex.EncodeToString(ctx.Sum(nil)))
}

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

// 生成32位大写MD5值
func MD5Upper(text string) string {
	result := MD5(text)
	return strings.ToUpper(result)
}

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

// 生成16进制的Sha1
func Sha1X(text string) string {
	ctx := sha1.New()
	_, err := ctx.Write([]byte(text))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x\n", ctx.Sum(nil))
}
