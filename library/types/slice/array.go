package slice

import (
	"github.com/gogo/protobuf/sortkeys"
	"sort"
)

//字符串元素是否在数组中
func InArray(find interface{}, arr interface{}) bool {
	if arr == nil {
		return false
	}
	switch find.(type) {
	case int:
		return InArrayInt(find.(int), arr.([]int))
	case int8:
		return InArrayInt8(find.(int8), arr.([]int8))
	case int16:
		return InArrayInt16(find.(int16), arr.([]int16))
	case int32:
		return InArrayInt32(find.(int32), arr.([]int32))
	case int64:
		return InArrayInt64(find.(int64), arr.([]int64))
	case uint:
		return InArrayUint(find.(uint), arr.([]uint))
	case uint8:
		return InArrayUint8(find.(uint8), arr.([]uint8))
	case uint16:
		return InArrayUint16(find.(uint16), arr.([]uint16))
	case uint32:
		return InArrayUint32(find.(uint32), arr.([]uint32))
	case uint64:
		return InArrayUint64(find.(uint64), arr.([]uint64))
	case float32:
		return InArrayFloat32(find.(float32), arr.([]float32))
	case float64:
		return InArrayFloat64(find.(float64), arr.([]float64))
	case string:
		return InArrayString(find.(string), arr.([]string))
	}

	return false
}

func InArrayInt(find int, values []int) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayInt8(find int8, values []int8) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayInt16(find int16, values []int16) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayInt32(find int32, values []int32) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayInt64(find int64, values []int64) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayUint(find uint, values []uint) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayUint8(find uint8, values []uint8) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayUint16(find uint16, values []uint16) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}
func InArrayUint32(find uint32, values []uint32) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayUint64(find uint64, values []uint64) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayFloat32(find float32, values []float32) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayFloat64(find float64, values []float64) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func InArrayString(find string, values []string) bool {
	for _, v := range values {
		if find == v {
			return true
		}
	}
	return false
}

func MaxInt64(values []int64) int64 {
	if values == nil || len(values) == 0 {
		return 0
	}
	sortkeys.Int64s(values)
	return values[len(values)-1]
}

func MinInt64(values []int64) int64 {
	if values == nil || len(values) == 0 {
		return 0
	}
	sortkeys.Int64s(values)
	return values[0]
}

func AvgInt64(values []int64) int64 {
	if values == nil || len(values) == 0 {
		return 0
	}
	var total int64
	for _, i := range values {
		total += i
	}
	return total / int64(len(values))
}

func SumInt64(values []int64) int64 {
	if values == nil || len(values) == 0 {
		return 0
	}
	var total int64
	for _, i := range values {
		total += i
	}
	return total
}

func MaxInt(values []int) int {
	if values == nil || len(values) == 0 {
		return 0
	}
	sort.Ints(values)
	return values[len(values)-1]
}

func MinInt(values []int) int {
	if values == nil || len(values) == 0 {
		return 0
	}
	sort.Ints(values)
	return values[0]
}

func AvgInt(values []int) int {
	if values == nil || len(values) == 0 {
		return 0
	}
	var total int
	for _, i := range values {
		total += i
	}
	return total / len(values)
}

func SumInt(values []int) int {
	if values == nil || len(values) == 0 {
		return 0
	}
	var total int
	for _, i := range values {
		total += i
	}
	return total
}

func MaxFloat64(values []float64) float64 {
	if values == nil || len(values) == 0 {
		return 0
	}
	sortkeys.Float64s(values)
	return values[len(values)-1]
}

func MinFloat64(values []float64) float64 {
	if values == nil || len(values) == 0 {
		return 0
	}
	sortkeys.Float64s(values)
	return values[0]
}

func AvgFloat64(values []float64) float64 {
	if values == nil || len(values) == 0 {
		return 0
	}
	var total float64
	for _, i := range values {
		total += i
	}
	return total / float64(len(values))
}

func SumFloat64(values []float64) float64 {
	if values == nil || len(values) == 0 {
		return 0
	}
	var total float64
	for _, i := range values {
		total += i
	}
	return total
}
