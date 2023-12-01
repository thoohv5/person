package util

import (
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

// UpperCamelName 大驼峰
func UpperCamelName(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// LowerCameName 小驼峰
func LowerCameName(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// Strikethrough 中划线
func Strikethrough(s string) string {
	return strings.ReplaceAll(s, "_", "-")
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

// Strikethrough2Underline 中划线转下划线
func Strikethrough2Underline(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}
