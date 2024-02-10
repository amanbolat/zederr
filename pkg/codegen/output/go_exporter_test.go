package output_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/amanbolat/zederr/pkg/codegen/core"
	"github.com/amanbolat/zederr/pkg/codegen/output"
	"github.com/amanbolat/zederr/pkg/codegen/parser"
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

	expectedGoErrorsFile, err := os.ReadFile("testdata/test01/errors.go")
	require.NoError(t, err)

	expectedGoEmbedFile, err := os.ReadFile("testdata/test01/embed.go")
	require.NoError(t, err)

	expectedLocaleEn, err := os.ReadFile("testdata/test01/locale.en.toml")
	require.NoError(t, err)

	expectedLocaleRu, err := os.ReadFile("testdata/test01/locale.ru.toml")
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	cfg := core.GoExporterConfig{
		PackageName: "zederrtest",
		Output:      buf,
		OutputPath:  "",
	}
	err = e.Export(cfg, zedErrs)
	require.NoError(t, err)

	fmt.Println(buf.String())

	expectedBuf := bytes.NewBuffer(nil)
	expectedBuf.Write(expectedGoErrorsFile)
	expectedBuf.Write(expectedGoEmbedFile)
	expectedBuf.Write(expectedLocaleEn)
	expectedBuf.Write(expectedLocaleRu)
	assert.Equal(t, expectedBuf.String(), buf.String())
}
