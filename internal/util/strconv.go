// Package util .
package util

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// MustParseInt .
func MustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("strconv.Atoi(%s) error: %v", s, err))
	}
	return i
}

// Int32SliceToStringSlice []int32 to []string
func Int32SliceToStringSlice(int32Slice []int32) []string {
	var stringSlice []string
	for _, v := range int32Slice {
		stringSlice = append(stringSlice, strconv.Itoa(int(v)))
	}
	return stringSlice
}

// HexToBin 16转2，不足8位补0。字符串长度必须为偶数。不是偶数的话，前面会补0。
func HexToBin(str string) (string, error) {
	rest := ""
	if len(str)%2 != 0 {
		str = "0" + str
	}
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}
	for _, b := range bytes {
		rest += fmt.Sprintf("%08b", b)
	}
	return rest, nil
}

// BinToHex 2转16，不足2位补0。字符串长度必须为8的倍数。不是8的倍数的话，前面会补0。
func BinToHex(str string) (string, error) {
	rest := ""
	for i := 0; i < len(str)%8; i++ {
		str = "0" + str
	}
	for i := 0; i < len(str); i += 8 {
		v, err := strconv.ParseInt(str[i:i+8], 2, 64)
		if err != nil {
			return "", err
		}
		rest += fmt.Sprintf("%02x", v)
	}
	return rest, nil
}

// FormatHex 根据长度格式化16进制字符串，不足补0，超出从右边截断
func FormatHex(str string, length int32) string {
	if str == "" {
		return ""
	}
	if len(str) > int(length) {
		return str[len(str)-int(length):]
	}
	tmpLen := int(length) - len(str)
	for i := 0; i < tmpLen; i++ {
		str = "0" + str
	}
	return str
}
