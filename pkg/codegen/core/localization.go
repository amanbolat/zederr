package core

import (
	"fmt"
	"maps"
	"strings"
	"unicode/utf8"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/language"
)

type Localization struct {
	description     map[language.Tag]string
	title           map[language.Tag]string
	arguments       map[string]map[language.Tag]string
	publicMessage   map[language.Tag]string
	internalMessage map[language.Tag]string
	deprecated      map[language.Tag]string
}

func (l Localization) Description() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(l.description, m)

	return m
}

func (l Localization) Title() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(l.title, m)

	return m
}

func (l Localization) Arguments() map[string]map[language.Tag]string {
	m := map[string]map[language.Tag]string{}
	for argName, translations := range l.arguments {
		m[argName] = map[language.Tag]string{}
		maps.Copy(translations, m[argName])
	}

	return m
}

func (l Localization) PublicMessage() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(l.publicMessage, m)

	return m
}

func (l Localization) InternalMessage() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(l.internalMessage, m)

	return m
}

func (l Localization) Deprecated() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(l.deprecated, m)

	return m
}

func NewLocalization() Localization {
	return Localization{
		description:     map[language.Tag]string{},
		title:           map[language.Tag]string{},
		arguments:       map[string]map[language.Tag]string{},
		publicMessage:   map[language.Tag]string{},
		internalMessage: map[language.Tag]string{},
		deprecated:      map[language.Tag]string{},
	}
}

func (l *Localization) AddDescriptionTranslation(lang string, val string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}

	if !utf8.ValidString(val) {
		return fmt.Errorf("description is not a valid UTF-8 string; got %s", val)
	}

	l.description[tag] = val

	return nil
}

func (l *Localization) AddTitleTranslation(lang string, val string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}

	if !utf8.ValidString(val) {
		return fmt.Errorf("title is not a valid UTF-8 string; got %s", val)
	}

	l.title[tag] = val

	return nil
}

func (l *Localization) AddArgumentTranslation(name, lang, val string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("argument name is empty")
	}

	if !utf8.ValidString(val) {
		return fmt.Errorf("argument value is not a valid UTF-8 string; got %s", val)
	}

	name = strcase.ToLowerCamel(name)

	if _, ok := l.arguments[name]; !ok {
		l.arguments[name] = map[language.Tag]string{}
	}

	l.arguments[name][tag] = val

	return nil
}

func (l *Localization) AddPublicMessageTranslation(lang string, val string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}

	if !utf8.ValidString(val) {
		return fmt.Errorf("public message is not a valid UTF-8 string; got %s", val)
	}

	l.publicMessage[tag] = val

	return nil
}

func (l *Localization) AddInternalMessageTranslation(lang string, val string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}

	if !utf8.ValidString(val) {
		return fmt.Errorf("internal message is not a valid UTF-8 string; got %s", val)
	}

	l.internalMessage[tag] = val

	return nil
}

func (l *Localization) AddDeprecatedTranslation(lang string, val string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}

	l.deprecated[tag] = val

	return nil
}
