package core

import (
	"golang.org/x/text/language"
)

type Spec struct {
	Version       string
	Domain        string
	Namespace     string
	DefaultLocale language.Tag
	Errors        []Error
}
