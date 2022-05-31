// Author: chenqionghe
// Time: 2018-10
// 字符串相关操作

package str

import (
	"app/library/types/convert"
	"net"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// 驼峰式写法转为下划线写法
func Camel2Case(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

// 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

func Ucfirst(str string) string {
	// 首字母大写
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

//TODO 获取搜索的多值,支持中文逗号/英文逗号/空格
func ParseMultiInt(value string) []int {
	return []int{convert.MustInt(value)}
}

//表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

//表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

//过滤非法字符串 去除空格、去除 TODO
func Filter(v string) string {
	return ""
}

//判断字符串是否为空
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

//是否是git地址
func IsIp(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil
}

//字符串截取，超过字符数显示...
func LimitTo(s string, limit int) string {
	runStr := []rune(s)
	length := len(runStr)
	if length < limit {
		return s
	}
	return string(runStr[:limit]) + "..."
}

//字符串截取，偏移
func OffsetTo(s string, offset int) string {
	runStr := []rune(s)
	length := len(runStr)
	if length < offset {
		return ""
	}
	return string(runStr[offset:])
}
