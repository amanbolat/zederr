package core

import (
	"fmt"
	"maps"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/language"
)

type Localization struct {
	description map[language.Tag]string
	arguments   map[string]map[language.Tag]string
	message     map[language.Tag]string
}

func NewLocalization() Localization {
	return Localization{
		description: map[language.Tag]string{},
		arguments:   map[string]map[language.Tag]string{},
		message:     map[language.Tag]string{},
	}
}

func (l Localization) AllLanguages() []language.Tag {
	arrSize := len(l.description) + len(l.message)

	for _, translations := range l.arguments {
		arrSize += len(translations)
	}

	tags := make([]language.Tag, 0, arrSize)

	for tag := range l.description {
		tags = append(tags, tag)
	}

	for tag := range l.message {
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

func (l Localization) Arguments() map[string]map[language.Tag]string {
	args := map[string]map[language.Tag]string{}

	for argName, translations := range l.arguments {
		args[argName] = map[language.Tag]string{}
		maps.Copy(args[argName], translations)
	}

	return args
}

func (l Localization) Message() map[language.Tag]string {
	m := map[language.Tag]string{}
	maps.Copy(m, l.message)

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

func (l *Localization) AddMessageTranslation(lang string, val string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}

	if !utf8.ValidString(val) {
		return fmt.Errorf("public message is not a valid UTF-8 string; got %s", val)
	}

	l.message[tag] = val

	return nil
}
