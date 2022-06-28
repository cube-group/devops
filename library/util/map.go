package util

import "strings"

//map[string]string转为字符串
func MapToString(m map[string]string) string {
	if m == nil || len(m) == 0 {
		return ""
	}

	s := make([]string, 0)
	for k, v := range m {
		s = append(s, k+"="+v)
	}
	return strings.Join(s, ",")
}

//string转为map[string]string
func StringToMap(s string) map[string]string {
	m := make(map[string]string, 0)
	if s == "" {
		return m
	}

	arr := strings.Split(s, ",")
	for _, v := range arr {
		arr2 := strings.Split(v, "=")
		m[arr2[0]] = arr2[1]
	}
	return m
}

func MapContact(option ...map[string]string) map[string]string {
	var res = make(map[string]string)
	for _, item := range option {
		for k, v := range item {
			res[k] = v
		}
	}
	return res
}
