package util

import (
	"strconv"
	"strings"
)

// SplitInt32 拆分
func SplitInt32(s, sep string) ([]int32, error) {
	split := strings.Split(s, sep)
	arr := make([]int32, 0, len(split))
	for _, item := range split {
		atoi, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		arr = append(arr, int32(atoi))
	}
	return arr, nil
}

// DuplicateRemoval 去重
func DuplicateRemoval(arr []int32) []int32 {
	check := make(map[int32]interface{})
	ret := make([]int32, 0, len(arr))

	for _, item := range arr {
		if _, ok := check[item]; !ok {
			ret = append(ret, item)
		}
	}
	return ret
}

// DiffInt32 取差集
func DiffInt32(arr1, arr2 []int32) []int32 {
	check := make(map[int32]interface{})
	ret := make([]int32, 0, len(arr1))

	for _, item := range arr2 {
		check[item] = nil
	}

	for _, item := range arr1 {
		if _, ok := check[item]; !ok {
			ret = append(ret, item)
		}
	}
	return ret
}

// SortInt32 排序
func SortInt32(arr []int32, desc bool) []int32 {
	if desc {
		for i := 0; i < len(arr); i++ {
			for j := i + 1; j < len(arr); j++ {
				if arr[i] < arr[j] {
					arr[i], arr[j] = arr[j], arr[i]
				}
			}
		}
	} else {
		for i := 0; i < len(arr); i++ {
			for j := i + 1; j < len(arr); j++ {
				if arr[i] > arr[j] {
					arr[i], arr[j] = arr[j], arr[i]
				}
			}
		}
	}
	return arr
}

// DiffInt32Plus 加强版取差集
func DiffInt32Plus(slice1, slice2 []int32) []int32 {
	// 创建一个 map 用于存储较长切片中的元素
	elements := make(map[int32]struct{})
	var result []int32

	// 将较长的切片中的元素添加到 map 中
	var longer, shorter []int32
	if len(slice1) >= len(slice2) {
		longer = slice1
		shorter = slice2
	} else {
		longer = slice2
		shorter = slice1
	}

	for _, num := range shorter {
		elements[num] = struct{}{}
	}

	// 检查较短切片中的元素是否在 map 中，如果不在则加入到结果中
	for _, num := range longer {
		if _, ok := elements[num]; !ok {
			result = append(result, num)
		}
	}

	return result
}

// InInt32ASlice 判断是否在[]int32中
func InInt32ASlice(arr []int32, item int32) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}
