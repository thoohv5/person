package util

import (
	"fmt"
	"strings"
)

const (
	_customKey = "key_"
)

// MakeCustomField 生成自定义字段
func MakeCustomField(key string) string {
	return fmt.Sprintf("%s%s", _customKey, key)
}

// IsCustomField 是否为自定义字段
func IsCustomField(key string) bool {
	return strings.HasPrefix(key, "key_")
}
