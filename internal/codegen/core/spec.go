package core

import (
	"golang.org/x/text/language"
)

type Spec struct {
	Version       string
	DefaultLocale language.Tag
	Errors        []Error
}

func (s Spec) HasTimestampArguments() bool {
	for _, err := range s.Errors {
		for _, arg := range err.Arguments() {
			if arg.Typ() == ArgumentTypeTimestamp {
				return true
			}
		}
	}

	return false
}
