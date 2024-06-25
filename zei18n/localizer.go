package zei18n

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/amanbolat/zederr/zeerr"
)

type localizer struct {
	localizers  map[language.Tag]*i18n.Localizer
	defaultLang language.Tag
}

// NewLocalizer creates a new localizer.
// NOTE: it's meant to be used only by the generated code.
func NewLocalizer(defaultLocale string, messagesMap map[string][]byte) (zeerr.Localizer, error) {
	defaultLocaleTag, err := language.Parse(defaultLocale)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default locale [%s]: %w", defaultLocale, err)
	}

	loc := localizer{
		defaultLang: defaultLocaleTag,
		localizers:  map[language.Tag]*i18n.Localizer{},
	}

	bundle := i18n.NewBundle(defaultLocaleTag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	var defaultLangFound bool

	for lang, data := range messagesMap {
		langTag, err := language.Parse(lang)
		if err != nil {
			return nil, fmt.Errorf("failed to parse language tag from string [%s]: %w", lang, err)
		}

		if langTag == defaultLocaleTag {
			defaultLangFound = true
		}

		bundlePath := fmt.Sprintf("%s.toml", lang)
		bundle.MustParseMessageFileBytes(data, bundlePath)

		loc.localizers[langTag] = i18n.NewLocalizer(bundle, langTag.String())
	}

	if !defaultLangFound {
		return nil, fmt.Errorf("bundle has no messages for default locale [%s]", defaultLocaleTag)
	}

	return &loc, nil
}

// LocalizeMessage localizes error's public message.
func (l *localizer) LocalizeMessage(id string, lang language.Tag, args map[string]any) string {
	loc, ok := l.localizers[lang]
	if !ok {
		loc = l.localizers[l.defaultLang]
	}

	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    id + "_message",
		TemplateData: args,
	})
	if err != nil {
		return ""
	}

	return msg
}
