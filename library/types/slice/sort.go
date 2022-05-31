package slice

import (
	"sort"
)

type sortMapSlice struct {
	key    string
	values []map[string]interface{}
}

func (s sortMapSlice) Len() int {
	return len(s.values)
}

func (s sortMapSlice) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

func (s sortMapSlice) Less(i, j int) bool {
	iV, iOk := s.values[i][s.key].(int)
	if !iOk {
		return false
	}
	jV, jOk := s.values[j][s.key].(int)
	if !jOk {
		return false
	}
	return iV < jV
}

func SortMapSliceByKeyInt(values []map[string]interface{}, key string) []map[string]interface{} {
	slice := sortMapSlice{key: key, values: values}
	sort.Stable(slice)
	return slice.values
}
