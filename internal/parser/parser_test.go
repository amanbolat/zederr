package parser_test

import (
	"strings"
	"testing"

	"github.com/amanbolat/zederr/internal/core"
	"github.com/amanbolat/zederr/internal/parser"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func TestParser(t *testing.T) {
	p := parser.NewParser("{{", "}}", &parser.Config{Debug: false})

	testCases := []struct {
		name                    string
		actualTranslation       string
		expecterLocalizableText string
		expectedFields          []core.Param
		shouldFail              bool
	}{
		{
			name:                    "single field with no type",
			actualTranslation:       "Error with {{ .Param1 }}.",
			expecterLocalizableText: "Error with {{.Param1}}.",
			expectedFields:          []core.Param{{Name: "Param1", Type: "any"}},
			shouldFail:              false,
		},
		{
			name:                    "wrong parameter name – dots are not allowed",
			actualTranslation:       "Error with {{ .Param1.Param2 }}.",
			expecterLocalizableText: "",
			expectedFields:          []core.Param{},
			shouldFail:              true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fieldMap, actualLocalizableText, err := p.Parse(testCase.actualTranslation)
			if testCase.shouldFail {
				assert.Error(t, err)
			}

			actualFields := maps.Values(fieldMap)

			slices.SortFunc(actualFields, func(i, j core.Param) int {
				return strings.Compare(i.Name, j.Name)
			})
			assert.Equal(t, testCase.expectedFields, actualFields)
			assert.Equal(t, testCase.expecterLocalizableText, actualLocalizableText)
		})
	}
}
