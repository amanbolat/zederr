package core

import (
	"fmt"
	"maps"
	"strings"
	"unicode/utf8"

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

func (l Localization) AllLanguages() []language.Tag {
	arrSize := len(l.description) + len(l.title) + len(l.publicMessage) + len(l.internalMessage) + len(l.deprecated)

	for _, translations := range l.arguments {
		arrSize += len(translations)
	}

	tags := make([]language.Tag, 0, arrSize)

	for tag := range l.description {
		tags = append(tags, tag)
	}

	for tag := range l.title {
		tags = append(tags, tag)
	}

	for tag := range l.publicMessage {
		tags = append(tags, tag)
	}

	for tag := range l.internalMessage {
		tags = append(tags, tag)
	}

	for tag := range l.deprecated {
		tags = append(tags, tag)
	}

	for _, translations := range l.arguments {
		for tag := range translations {
			tags = append(tags, tag)
		}
	}

	return tags
}

func (l Localization) Description() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(m, l.description)

	return m
}

func (l Localization) Title() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(m, l.title)

	return m
}

func (l Localization) Arguments() map[string]map[language.Tag]string {
	args := map[string]map[language.Tag]string{}

	for argName, translations := range l.arguments {
		args[argName] = map[language.Tag]string{}
		maps.Copy(args[argName], translations)
	}

	return args
}

func (l Localization) PublicMessage() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(m, l.publicMessage)

	return m
}

func (l Localization) InternalMessage() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(m, l.internalMessage)

	return m
}

func (l Localization) Deprecated() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(m, l.deprecated)

	return m
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

func (l *Localization) AddArgumentTranslation(name, description, lang string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return fmt.Errorf("failed to parse language (%s) tag; %w", lang, err)
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("argument name is empty")
	}

	if !utf8.ValidString(description) {
		return fmt.Errorf("argument description is not a valid UTF-8 string; got %s", description)
	}

	if _, ok := l.arguments[name]; !ok {
		l.arguments[name] = map[language.Tag]string{}
	}

	l.arguments[name][tag] = description

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
