package zei18n

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/amanbolat/zederr/zeerr"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type localizer struct {
	localizers  map[language.Tag]*i18n.Localizer
	defaultLang language.Tag
}

func NewLocalizer(defaultLocale string, messagesMap map[string][]byte) (zeerr.Localizer, error) {
	defaultLocaleTag, err := language.Parse(defaultLocale)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default locale [%s]: %w", defaultLocale, err)
	}

	l := localizer{
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

		l.localizers[langTag] = i18n.NewLocalizer(bundle, langTag.String())
	}

	if !defaultLangFound {
		return nil, fmt.Errorf("bundle has no messages for default locale [%s]", defaultLocaleTag)
	}

	return &l, nil
}

func (l *localizer) LocalizePublicMessage(errUID string, lang language.Tag, args zeerr.Arguments) string {
	loc, ok := l.localizers[lang]
	if !ok {
		loc = l.localizers[l.defaultLang]
	}

	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    errUID + "_public_msg",
		TemplateData: args,
	})
	if err != nil {
		return ""
	}

	return msg
}

func (l *localizer) LocalizeError(e zeerr.Error, lang language.Tag) (string, bool, error) {
	loc, ok := l.localizers[lang]
	if !ok {
		loc = l.localizers[l.defaultLang]
	}

	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    e.UID(),
		TemplateData: e.Args(),
	})
	if err != nil {
		return "", false, fmt.Errorf("failed to localize error [%s]: %w", e.UID(), err)
	}

	return msg, ok, nil
}
