// Author: chenqionghe
// Time: 2018-10
// 自定义模板函数

package ginutil

import (
	"app/library/crypt/base64"
	"app/library/types/convert"
	"app/library/types/str"
	"app/library/types/times"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"reflect"
	"time"
)

//模板函数
func ViewFunc() template.FuncMap {
	return template.FuncMap{
		"datetime":      times.FormatDatetime,
		"static":        Static,
		"safe":          Safe,
		"timeFormat":    times.TimeFormat,
		"strtotime":     times.Strtotime,
		"strTimeFormat": times.StrTimeFormat,
		"timestampToTime": func(v interface{}) time.Time {
			return time.Unix(convert.MustInt64(v), 0)
		},
		"checked":    Checked,
		"selected":   Selected,
		"isSelected": IsSelected,
		"available":  Available,
		"disabled":   Disabled,

		"UnicodeEmojiDecode": str.UnicodeEmojiCode,
		"UnicodeEmojiCode":   str.UnicodeEmojiCode,
		"SizeFormat":         convert.SizeFormat,
		"MBSizeFormat":       convert.MBSizeFormat,
		"IsStrEmpty":         str.IsEmpty,
		"lastIndex":          func(size int) int { return size - 1 },
		"Iterate": func(count uint) []uint {
			var i uint
			var Items []uint
			for i = 1; i <= (count); i++ {
				Items = append(Items, i)
			}
			return Items
		},
		"toStr": convert.MustString,
		"isLoghubAttr": func(v string) bool {
			length := len(v)
			if length < 4 {
				return false
			}
			if v[:2] == "__" && v[length-2:] == "__" {
				return true
			}
			return false
		},

		"limitTo":      str.LimitTo,
		"offsetTo":     str.OffsetTo,
		"base64Encode": base64.Base64Encode,
		"base64Dncode": base64.Base64Decode,
	}
}

//获取公共文件
func Static(path string) string {
	return "/public" + path + "?" + viper.GetString("version")
}

//不转义输出
func Safe(x string) interface{} {
	return template.HTML(x)
}

//是否选中select
func Selected(a interface{}, b interface{}) string {
	if IsSelected(a, b) {
		return "selected"
	}
	return ""
}

//是否选中select
func Checked(a interface{}, b interface{}) string {
	if fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b) {
		return "checked"
	}
	return ""
}

//是否选中select
func Disabled(data interface{}, name string) string {
	if keyExists(data, name) {
		return "disabled"
	}
	return ""
}

//是否选中select
func IsSelected(a interface{}, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func Available(data interface{}, name string) bool {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	return v.FieldByName(name).IsValid()
}

func keyExists(data interface{}, name string) bool {
	dataMap, ok := data.(map[string]interface{})
	if ok {
		if _, ok := dataMap[name]; ok {
			return true
		}
	}
	return false
}
