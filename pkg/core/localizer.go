package core

import (
	"golang.org/x/text/language"
)

// Localizer is responsible for localizing public and internal error messages.
// It should only bundle the locales for a single namespace.
type Localizer interface {
	LocalizePublicMessage(errUID string, lang language.Tag, args Arguments) string
}
