package output_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/amanbolat/zederr/internal/core"
	"github.com/amanbolat/zederr/internal/output"
	"github.com/amanbolat/zederr/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoExporter(t *testing.T) {
	e := output.GoExporter{}

	p := parser.NewParser("{{", "}}", nil)
	b := core.NewErrorBuilder(p)

	zedErr, err := b.NewError(
		"core.error",
		2,
		400,
		"core error",
		map[string]string{
			"en": "core error {{ .Error | string }}",
			"ru": "ошибка ядра {{ string .User }}",
		})
	require.NoError(t, err)
	zedErrs := []core.Error{
		zedErr,
	}

	expectedFile, err := os.ReadFile("testdata/test_output_1.go")
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	cfg := core.GoExporterConfig{
		PackageName: "zederrtest",
		Output:      buf,
		OutputPath:  "",
	}
	err = e.Export(cfg, zedErrs)
	require.NoError(t, err)

	assert.Equal(t, string(expectedFile), buf.String())
}
