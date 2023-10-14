package localize

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	ic "github.com/thoohv5/person/internal/context"
)

type (
	localize struct {
		label   string
		printer *message.Printer
	}
	ILocalize interface {
		Translate(ctx context.Context, key message.Reference, args ...interface{}) string
	}
)

var (
	locales = []*localize{
		{
			label:   ZH,
			printer: message.NewPrinter(language.MustParse(ZH)),
		},
		{
			label:   EN,
			printer: message.NewPrinter(language.MustParse(EN)),
		},
	}
	defaultLocale = MustGet(ZH)
)

const (
	ZH = "zh-CN"
	EN = "en-GB"
)

func MustGet(label string) ILocalize {
	for _, locale := range locales {
		if label == locale.label {
			return locale
		}
	}
	return locales[len(locales)-1]
}

func Get(label string) (ILocalize, bool) {
	for _, locale := range locales {
		if label == locale.label {
			return locale, true
		}
	}
	return nil, false
}

func (l *localize) Translate(ctx context.Context, key message.Reference, args ...interface{}) string {
	locale := l

	if gtx, ok := ctx.(*gin.Context); ok {
		ctx = gtx.Request.Context()
	}

	if loc, ok := Get(ic.FromCtxLang(ctx)); ok {
		locale = loc.(*localize)
	}
	return locale.printer.Sprintf(key, args...)
}

func Translate(ctx context.Context, key message.Reference, args ...interface{}) string {
	return defaultLocale.Translate(ctx, key, args...)
}

func ChangeLocale(label string) (ILocalize, bool) {
	dl, ok := Get(label)
	if ok {
		defaultLocale = dl
		return defaultLocale, true
	}
	return defaultLocale, false
}
