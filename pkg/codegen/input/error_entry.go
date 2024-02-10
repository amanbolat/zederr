package input

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"gopkg.in/yaml.v3"
)

// ErrorEntry represents a single error entry in the error codes file.
// It is used only for unmarshalling from the source file.
type ErrorEntry struct {
	ID           string            `yaml:"id"`
	GRPCCode     codes.Code        `yaml:"grpc_code"`
	HTTPCode     int               `yaml:"http_code"`
	Description  string            `yaml:"description"`
	Translations map[string]string `yaml:"translations"`
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

		if err := value.Content[i].Decode(&entry.ID); err != nil {
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
