package util

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// IString .
type IString string

// InitialUpper 首字母转大写
func InitialUpper(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(string(s[0])) + s[1:]
	}
	return s
}

// InitialLower 首字母转小写
func InitialLower(s string) string {
	if len(s) > 0 {
		return strings.ToLower(string(s[0])) + s[1:]
	}
	return s
}

// 特殊字段
var special = map[string]string{
	"ID":  "Id",
	"IP":  "Ip",
	"OUI": "Oui",
}

// Camel2Underline 驼峰转下划线
func Camel2Underline(s string) string {
	var result []rune
	for k, v := range special {
		s = strings.ReplaceAll(s, k, v)
	}
	for i, v := range s {
		if i > 0 && v >= 'A' && v <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, v)
	}
	return strings.ToLower(string(result))
}

// Camel2Strikethrough 驼峰转中划线
func Camel2Strikethrough(s string) string {
	var result []rune
	for i, v := range s {
		if i > 0 && v >= 'A' && v <= 'Z' {
			result = append(result, '-')
		}
		result = append(result, v)
	}
	return strings.ToLower(string(result))
}

// Strikethrough2Underline 中划线转下划线
func Strikethrough2Underline(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

// RangeExtractNumberString 将一段数字字符串转换成简短的描述
// example: 1,2,3,4,7  => 1-4, 7
func RangeExtractNumberString(original string) string {
	strArr := strings.Split(original, ",")
	intArr := []int{}
	for _, str := range strArr {
		intVal, err := strconv.Atoi(str)
		if err != nil {
			continue
		}
		intArr = append(intArr, intVal)
	}

	sort.Ints(intArr)
	return extractNumberRange(intArr)
}

func extractNumberRange(list []int) string {
	if len(list) == 0 {
		return ""
	}

	result := []int{list[0]}
	tag := list[0] - 1 // 用一个数字来做占位分割符
	for i := 1; i < len(list); i++ {
		// 是连续的
		if list[i] == list[i-1]+1 {
			if len(result) > 1 && result[len(result)-2] == tag {
				result[len(result)-1] = list[i]
			} else {
				result = append(result, tag, list[i])
			}
		} else {
			// 非连续的
			result = append(result, list[i])
		}
	}

	// 组装输出
	output := ""
	for i := 0; i < len(result); i++ {
		if i < len(result)-2 && result[i+1] == tag {
			if result[i]+1 == result[i+2] {
				output += fmt.Sprintf("%d,", result[i])
			} else {
				output += fmt.Sprintf("%d-", result[i])
			}
			i++
		} else {
			if i == len(result)-1 {
				output += fmt.Sprintf("%d", result[i])
			} else {
				output += fmt.Sprintf("%d,", result[i])
			}
		}
	}
	return output
}
