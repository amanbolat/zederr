package input

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"gopkg.in/yaml.v3"
)

type ErrorListSpecification struct {
	SpecVersion   string       `yaml:"spec_version"`
	Domain        string       `yaml:"domain"`
	Namespace     string       `yaml:"namespace"`
	DefaultLocale string       `yaml:"default_locale"`
	Errors        ErrorEntries `yaml:"errors"`
}

type Argument struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
}

type Arguments []Argument

func (a *Arguments) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`arguments` should be of type yaml.MappingNode, but got %v", value.Kind)
	}

	*a = make([]Argument, len(value.Content)/2)
	for i := 0; i < len(value.Content); i += 2 {
		var entry = &(*a)[i/2]
		if err := value.Content[i+1].Decode(&entry); err != nil {
			return fmt.Errorf("failed to decode argument content: %w", err)
		}

		if err := value.Content[i].Decode(&entry.Name); err != nil {
			return fmt.Errorf("failed to decode argument name: %w", err)
		}
	}
	return nil
}

type Translation struct {
	Lang  string
	Value string
}

type Translations []Translation

func (t *Translations) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`translations` should be of type yaml.MappingNode, but got %v", value.Kind)
	}

	*t = make([]Translation, len(value.Content)/2)
	for i := 0; i < len(value.Content); i += 2 {
		var entry = &(*t)[i/2]
		if err := value.Content[i+1].Decode(&entry.Value); err != nil {
			return err
		}

		if err := value.Content[i].Decode(&entry.Lang); err != nil {
			return err
		}
	}
	return nil

}

type LocalizationArgument struct {
	Name        string       `yaml:"name"`
	Description Translations `yaml:"description"`
}

type LocalizationArguments []LocalizationArgument

func (a *LocalizationArguments) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`arguments` should be of type yaml.MappingNode, but got %v", value.Kind)
	}

	*a = make([]LocalizationArgument, len(value.Content)/2)
	for i := 0; i < len(value.Content); i += 2 {
		var entry = &(*a)[i/2]
		if err := value.Content[i+1].Decode(&entry); err != nil {
			return err
		}

		if err := value.Content[i].Decode(&entry.Name); err != nil {
			return err
		}
	}

	return nil
}

type Localization struct {
	Arguments       LocalizationArguments `yaml:"arguments"`
	Description     Translations          `yaml:"description"`
	Title           Translations          `yaml:"title"`
	InternalMessage Translations          `yaml:"internal_message"`
	PublicMessage   Translations          `yaml:"public_message"`
	Deprecated      Translations          `yaml:"deprecated"`
}

// ErrorEntry represents a single error entry in the error codes file.
// It is used only for unmarshalling from the source file.
type ErrorEntry struct {
	Code            string        `yaml:"code"`
	GRPCCode        codes.Code    `yaml:"grpc_code"`
	HTTPCode        int           `yaml:"http_code"`
	Description     string        `yaml:"description"`
	Deprecated      string        `yaml:"deprecated"`
	Arguments       Arguments     `yaml:"arguments"`
	Title           string        `yaml:"title"`
	InternalMessage string        `yaml:"internal_message"`
	PublicMessage   string        `yaml:"public_message"`
	Localization    *Localization `yaml:"localization"`
}

// ErrorEntries is used to customize YAML unmarshalling of ErrorEntry.
type ErrorEntries []ErrorEntry

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (p *ErrorEntries) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`errors` should be of type yaml.MappingNode, but got %v", value.Kind)
	}

	*p = make([]ErrorEntry, len(value.Content)/2)
	for i := 0; i < len(value.Content); i += 2 {
		var entry = &(*p)[i/2]
		if err := value.Content[i+1].Decode(&entry); err != nil {
			return err
		}

		if err := value.Content[i].Decode(&entry.Code); err != nil {
			return err
		}
	}
	return nil
}

// AssertErrorEntriesEquality asserts that two ErrorEntries are equal.
// Used for testing only.
func AssertErrorEntriesEquality(t *testing.T, expected, actual ErrorEntries) {
	if assert.Len(t, actual, len(expected)) {
		for i := range expected {
			assert.Emptyf(t, cmp.Diff(expected[i], actual[i]), "expected and actual error entries are not equal")
		}
	}
}
