package validate

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type base struct {
	label []string
	text  string
}

var mFieldDoc = make(map[string]string)

// RegisterFieldDoc 注册
func RegisterFieldDoc(fd map[string]string) {
	for key, value := range fd {
		mFieldDoc[key] = value
	}
}

func (l *base) Tag() []string {
	return l.label
}

func (l *base) RegisterTranslationsFunc() validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		for _, tag := range l.Tag() {
			if err := ut.Add(tag, l.text, false); err != nil {
				return err
			}
		}
		return nil
	}
}

func (l *base) TranslationFunc() validator.TranslationFunc {
	return func(ut ut.Translator, fe validator.FieldError) string {
		for _, tag := range l.Tag() {
			ff := fe.Field()
			if fd, ok := mFieldDoc[fe.StructNamespace()]; ok {
				ff = fd
			}
			t, _ := ut.T(tag, ff, fe.Param())
			return t
		}
		return ""
	}
}
