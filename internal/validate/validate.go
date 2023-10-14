package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ten "github.com/go-playground/validator/v10/translations/en"
	tzh "github.com/go-playground/validator/v10/translations/zh"
)

type IValidate interface {
	Tag() []string
	Func() validator.Func
	RegisterTranslationsFunc() validator.RegisterTranslationsFunc
	TranslationFunc() validator.TranslationFunc
}

var registeredValidate = make([]IValidate, 0)
var currentTrans ut.Translator

func RegisterValidate(i IValidate) {
	registeredValidate = append(registeredValidate, i)
}

func InitValidate(locale string) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhn := zh.New()
		uni := ut.New(zhn, zhn, en.New())
		// this is usually know or extracted from http 'Accept-Language' header
		// also see uni.FindTranslator(...)
		trans, tok := uni.GetTranslator(locale)
		if !tok {
			return fmt.Errorf("%s not Translator", locale)
		}
		currentTrans = trans

		switch locale {
		case "en":
			if err := ten.RegisterDefaultTranslations(v, trans); err != nil {
				return err
			}
		default: // zh
			if err := tzh.RegisterDefaultTranslations(v, trans); err != nil {
				return err
			}

		}

		for _, validate := range registeredValidate {
			for _, tag := range validate.Tag() {
				err := v.RegisterValidation(tag, validate.Func())
				if err != nil {
					return nil
				}
				if err = v.RegisterTranslation(tag, trans, validate.RegisterTranslationsFunc(), validate.TranslationFunc()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Translate(err error) (ret interface{}) {
	var errs validator.ValidationErrors
	switch {
	case errors.As(err, &errs):
		ret = removeTopStruct(errs.Translate(currentTrans))
	default:
		ret = err.Error()
	}
	return ret
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
