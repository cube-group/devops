// Author: chenqionghe
// Time: 2018-10
// json编码

package ginutil

import (
	"encoding/json"
)

//输出漂亮的json
func PrettyJson(v string) (string, error) {
	var res map[string]interface{}
	if err := json.Unmarshal([]byte(v), &res); err != nil {
		return "", err
	}

	resBytes, err := json.MarshalIndent(res, "", "    ")
	return string(resBytes), err
}

// json转换
// 将struct/map/slice转换为json字符串
// @param target interface{} 为struct/map/slice实例
// @return string
func JsonEncode(v interface{}) (string, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
