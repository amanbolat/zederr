package zeerr

import (
	"context"

	"golang.org/x/text/language"
)

type LocaleCtxKeyType struct{}

func ContextWithLocale(ctx context.Context, lang language.Tag) context.Context {
	return context.WithValue(ctx, LocaleCtxKeyType{}, lang)
}

// Localizer is responsible for localizing public and internal error messages.
type Localizer interface {
	// LocalizePublicMessage localizes error's public message.
	// Usually it's used only during the error construction.
	LocalizePublicMessage(errUID string, lang language.Tag, args Arguments) string
}
