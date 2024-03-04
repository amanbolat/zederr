package output_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/amanbolat/zederr/internal/codegen/core"
	"github.com/amanbolat/zederr/internal/codegen/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
)

func TestGoExporter(t *testing.T) {
	e := output.GoExporter{}

	b, err := core.NewErrorBuilder("1", "acme.com", "ns", "en")
	require.NoError(t, err)

	localization := core.NewLocalization()

	arg1, err := core.NewArgument("arg_1", "description", "string")
	require.NoError(t, err)

	zedErr, err := b.NewError(
		"code",
		codes.Canceled,
		400,
		"description",
		"title",
		"public message {{ .arg_1 }}",
		"internal message {{ .arg_1 }}",
		"",
		[]core.Argument{
			arg1,
		},
		localization,
	)
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
	cfg := core.ExportGo{
		PackageName: "zederrtest",
		Output:      buf,
		OutputPath:  "",
	}
	spec := core.Spec{
		Version:       "1",
		Domain:        "acme.com",
		Namespace:     "namespace",
		DefaultLocale: language.English,
		Errors:        zedErrs,
	}

	err = e.Export(cfg, spec)
	require.NoError(t, err)

	fmt.Println(buf.String())

	expectedBuf := bytes.NewBuffer(nil)
	expectedBuf.Write(expectedGoErrorsFile)
	expectedBuf.Write(expectedGoEmbedFile)
	expectedBuf.Write(expectedLocaleEn)
	expectedBuf.Write(expectedLocaleRu)
	assert.Equal(t, expectedBuf.String(), buf.String())
}
