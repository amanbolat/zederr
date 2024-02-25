package core_test

import (
	"testing"

	"github.com/amanbolat/zederr/internal/codegen/core"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	p := core.NewTemplateValidator(&core.TemplateValidatorConfig{
		Debug:     false,
		Arguments: map[string]struct{}{"Param": {}},
	})

	err := p.Validate("{{.Param}}")
	assert.NoError(t, err)
}
