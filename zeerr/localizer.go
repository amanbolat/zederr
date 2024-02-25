package zeerr

import (
	"context"

	"golang.org/x/text/language"
)

type LocaleCtxKeyType struct{}

func ContextWithLocale(ctx context.Context, lang language.Tag) context.Context {
	return context.WithValue(ctx, LocaleCtxKeyType{}, lang)
}

type Localizer interface {
	LocalizePublicMessage(errUID string, lang language.Tag, args Arguments) string
	LocalizeError(e Error, lang language.Tag) (string, bool, error)
}
