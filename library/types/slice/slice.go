package slice

import (
	"fmt"
	"sort"
	"strings"
)

//判断字符串是否存在于sice中
func SliceStringIn(s []string, find string) bool {
	if s == nil {
		return false
	}
	for _, v := range s {
		if v == find {
			return true
		}
	}
	return false
}

//字符串slice获取交集
func SliceStringMixed(s1, s2 []string) []string {
	if s1 == nil && s2 == nil {
		return nil
	} else if s1 == nil {
		return s2
	} else if s2 == nil {
		return s1
	}
	mixed := make([]string, 0)
	for _, v := range s1 {
		if SliceStringIn(s2, v) {
			mixed = append(mixed, v)
		}
	}
	sort.Strings(mixed)
	return mixed
}

//字符串slice获取并集
func SliceStringUnion(s1, s2 []string) []string {
	if s1 == nil && s2 == nil {
		return nil
	} else if s1 == nil {
		return s2
	} else if s2 == nil {
		return s1
	}
	union := make(map[string]bool)
	res := make([]string, 0)
	for _, v := range s1 {
		union[v] = true
	}
	for _, v := range s2 {
		union[v] = true
	}
	for k, _ := range union {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

//获取s1和s2的差集
//含义：在s1不在s2中的数组
func SliceStringDiff(s1, s2 []string) []string {
	if s1 == nil && s2 == nil {
		return nil
	} else if s1 == nil {
		return s2
	} else if s2 == nil {
		return s1
	}
	diff := make([]string, 0)
	for _, v := range s1 {
		if !SliceStringIn(s2, v) {
			diff = append(diff, v)
		}
	}
	sort.Strings(diff)
	return diff
}

//对比两个数组值是否完全相等
//包括顺序
func SliceStringEqual(s1, s2 []string) bool {
	if s1 == nil && s2 == nil {
		return true
	} else if s1 == nil || s2 == nil {
		return false
	}
	return strings.Join(s1, ",") == strings.Join(s2, ",")
}

/**
 * 删除某元素
 */
func SliceRemove(s []interface{}, item interface{}) []interface{} {
	for k := 0; k < len(s); k++ {
		if fmt.Sprintf("%v", s[k]) == fmt.Sprintf("%v", item) {
			s = append(s[:k], s[k+1:]...)
			k--
		}
	}
	return s
}

/************************************** 数组去重函数 **************************************/

func UniqueStringSlice(list []string) []string {
	res := make([]string, 0)
	keys := make(map[string]bool)
	for _, v := range list {
		if _, ok := keys[v]; ok {
			continue
		}
		keys[v] = true
		res = append(res, v)
	}
	return res
}

func UniqueIntSlice(list []int) []int {
	res := make([]int, 0)
	keys := make(map[int]bool)
	for _, v := range list {
		if _, ok := keys[v]; ok {
			continue
		}
		keys[v] = true
		res = append(res, v)
	}
	return res
}

func UniqueInt64Slice(list []int64) []int64 {
	res := make([]int64, 0)
	keys := make(map[int64]bool)
	for _, v := range list {
		if _, ok := keys[v]; ok {
			continue
		}
		keys[v] = true
		res = append(res, v)
	}
	return res
}

func UniqueUintSlice(list []uint) []uint {
	res := make([]uint, 0)
	keys := make(map[uint]bool)
	for _, v := range list {
		if _, ok := keys[v]; ok {
			continue
		}
		keys[v] = true
		res = append(res, v)
	}
	return res
}