package core

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	argumentNameRegex = regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$")
)

// Argument represents an argument used in the error messages.
type Argument struct {
	name        string
	description string
	typ         ArgumentType
}

func NewArgument(name, description, typ string) (Argument, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return Argument{}, fmt.Errorf("argument name is empty")
	}

	if !utf8.ValidString(name) {
		return Argument{}, fmt.Errorf("argument name is not a valid UTF-8 string; got %s", name)
	}

	if !argumentNameRegex.MatchString(name) {
		return Argument{}, fmt.Errorf("argument name is not a valid identifier; it should match regex pattern: %s; got %s", argumentNameRegex, name)
	}

	description = strings.TrimSpace(description)

	if description == "" {
		return Argument{}, fmt.Errorf("argument description is empty")
	}

	if !utf8.ValidString(description) {
		return Argument{}, fmt.Errorf("argument description is not a valid UTF-8 string; got %s", description)
	}

	argTyp, err := ParseArgumentType(typ)
	if err != nil {
		return Argument{}, err
	}

	return Argument{
		name:        name,
		description: description,
		typ:         argTyp,
	}, nil
}

func (a Argument) Name() string {
	return a.name
}

func (a Argument) Description() string {
	return a.description
}

func (a Argument) Typ() ArgumentType {
	return a.typ
}
