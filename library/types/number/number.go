package number

import (
	"fmt"
	"strconv"
)

//保留小数
//value 数字
//decimal 保留小数位数
func NumberDecimal(value interface{}, decimal ...interface{}) float64 {
	v := "2"
	if len(decimal) > 0 {
		v = fmt.Sprintf("%v", decimal[0])
	}
	res, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	if err != nil {
		fmt.Println("NumberDecimal1", value, v, err)
		return 0
	}
	res, err = strconv.ParseFloat(fmt.Sprintf("%."+v+"f", res), 64)
	if err != nil {
		fmt.Println("NumberDecimal2", value, v, err)
		return 0
	}
	return res
}

//保留两位小数并且返回字符串
func NumberDecimalToString(value float64) string {
	return fmt.Sprintf("%v", NumberDecimal(value))
}

func ToUint(value interface{}) uint {
	res, _ := strconv.Atoi(fmt.Sprintf("%v", value))
	return uint(res)
}

func ToInt(value interface{}) int {
	res, _ := strconv.Atoi(fmt.Sprintf("%v", value))
	return res
}

func ToFloat32(value interface{}) float32 {
	res, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 32)
	return float32(res)
}

func ToFloat64(value interface{}) float64 {
	res, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	return res
}

func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
