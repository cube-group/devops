package g

import (
	"app/library/consts"
	"app/library/crypt/base64"
	"app/library/types/convert"
	"app/library/types/str"
	"app/library/types/times"
	"fmt"
	"html/template"
	"reflect"
	"time"
)

//模板函数
func ViewFunc() template.FuncMap {
	return template.FuncMap{
		"datetime":   times.FormatDatetime,
		"safe":       Safe,
		"timeFormat": times.TimeFormat,
		"checked":    Checked,
		"selected":   Selected,
		"isSelected": IsSelected,
		"available":  Available,
		"disabled":   Disabled,

		//构建状态
		//"IsStatusBuildFail": consts.IsStatusBuildFail,
		//"IsStatusBuildOk":   consts.IsStatusBuildOk,
		//"IsStatusBuilding":  consts.IsStatusBuilding,
		//"IsStatusBuildDone": consts.IsStatusBuildDone,
		//
		////部署状态
		//"IsStatusDeploying":  consts.IsStatusDeploying,
		//"IsStatusDeployFail": consts.IsStatusDeployFail,
		//"IsStatusDeployOk":   consts.IsStatusDeployOk,
		//
		////上线总状态
		//"IsStatusDone": consts.IsStatusDone,
		//
		////首页提示状态
		//"IsStatusOnline":            consts.IsStatusOnline,
		//"IsStatusWarning":           consts.IsStatusWarning,
		//"IsStatusSuccess":           consts.IsStatusSuccess,
		//"IsStatusInBuild":           consts.IsStatusInBuild,
		//"IsStatusInDeploy":          consts.IsStatusInDeploy,
		//"ProjectStatusCn":           consts.ProjectStatusCn,
		//"OnlineCn":                  consts.ProjectOnlineCn,
		"AccessLevelCn":             consts.AccessLevelCn,
		"AccessLevelPermission":     consts.AccessLevelPermission,

		"UnicodeEmojiDecode": str.UnicodeEmojiDecode,
		"UnicodeEmojiCode":   str.UnicodeEmojiCode,
		"SizeFormat":         convert.SizeFormat,
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
		"timestampToTime": func(v interface{}) time.Time {
			return time.Unix(convert.MustInt64(v), 0)
		},
		"strtotime":     times.Strtotime,
		"strTimeFormat": times.StrTimeFormat,
		"MBSizeFormat":  convert.MBSizeFormat,
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

//TODO
func keyExists(data interface{}, name string) bool {
	dataMap, ok := data.(map[string]interface{})
	if ok {
		if _, ok := dataMap[name]; ok {
			return true
		}
	}
	return false
}
