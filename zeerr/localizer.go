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
	// LocalizeMessage localizes error's message.
	LocalizeMessage(id string, lang language.Tag, args map[string]any) string
}
