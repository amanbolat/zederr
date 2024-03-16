package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amanbolat/zederr/internal/codegen/core"
)

func TestParser(t *testing.T) {
	p := core.NewTemplateValidator(&core.TemplateValidatorConfig{
		Debug:     false,
		Arguments: map[string]struct{}{"Param": {}},
	})

	err := p.Validate("{{.Param}}")
	assert.NoError(t, err)
}
