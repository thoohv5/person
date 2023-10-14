package validate

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

// telephone 电话号码
// example: telephone=1-32
type telephone struct {
	base
	compile *regexp.Regexp
}

func init() {
	RegisterValidate(&telephone{
		base: base{
			label: []string{"telephone"},
			text:  "请输入合法电话号码",
		},
		compile: regexp.MustCompile("^[^-][0-9-]*$"),
	})
}

func (l *telephone) Func() validator.Func {
	return func(fl validator.FieldLevel) bool {
		str, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		if str == "" {
			return true
		}
		if !l.compile.MatchString(str) {
			return false
		}
		i := utf8.RuneCountInString(str)
		params := strings.Split(fl.Param(), "-")
		if len(params) != 2 {
			return false
		}

		minLen, err := strconv.Atoi(params[0])
		if err != nil {
			return false
		}

		maxLen, err := strconv.Atoi(params[1])
		if err != nil {
			return false
		}
		return maxLen >= i && minLen <= i
	}
}
