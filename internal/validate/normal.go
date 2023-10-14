package validate

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type normal struct {
	base
	compile *regexp.Regexp
}

func init() {
	RegisterValidate(&normal{
		base: base{
			label: []string{"normal"},
			text:  "{0}请输入以非下划线开头的字母，汉字，数字，下划线的任意组合",
		},
		compile: regexp.MustCompile("^[a-zA-Z\u4E00-\u9FA5\uF900-\uFA2D][0-9a-zA-Z_\u4E00-\u9FA5\uF900-\uFA2D]*$"),
	})
}

func (l *normal) Func() validator.Func {
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
		return true
	}
}
