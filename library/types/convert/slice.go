// Author: chenqionghe
// Time: 2018-10
// slice转换相关操作

package convert

func UniqueStrSlice(array []string) []string {
	result := make([]string, 0, len(array))
	temp := map[string]struct{}{}
	for _, item := range array {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func UniqueIntSlice(array []int) []int {
	result := make([]int, 0, len(array))
	temp := map[int]struct{}{}
	for _, item := range array {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func StringToIntSlice(ids []string) []int {
	var intArray []int
	for _, id := range ids {
		intArray = append(intArray, MustInt(id))
	}
	return intArray
}

//转换成唯一的，并且没有空值的数组,用于where in 查询走索引
func IDArray(array []string) []uint {
	result := make([]uint, 0, len(array))
	temp := map[string]struct{}{}
	for _, item := range array {
		if _, ok := temp[item]; ok {
			continue
		}
		value := MustUint(item)
		if value == 0 {
			continue
		}
		temp[item] = struct{}{}
		result = append(result, value)
	}
	return result
}
