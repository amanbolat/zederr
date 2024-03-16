package output_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"

	"github.com/amanbolat/zederr/internal/codegen/core"
	"github.com/amanbolat/zederr/internal/codegen/output"
)

func TestGoExporter(t *testing.T) {
	exporter := output.GoExporter{}

	builder, err := core.NewErrorBuilder("1", "acme.com", "ns", "en")
	require.NoError(t, err)

	localization := core.NewLocalization()

	arg1, err := core.NewArgument("arg_1", "description", "string")
	require.NoError(t, err)

	zedErr, err := builder.NewError(
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

	tempFile, err := os.CreateTemp("", "zederrtest.*.txt")
	require.NoError(t, err)

	defer tempFile.Close()

	tempFileStat, err := tempFile.Stat()
	require.NoError(t, err)

	cfg := core.ExportGo{
		PackageName: "zederrtest",
		OutputPath:  tempFileStat.Name(),
	}
	spec := core.Spec{
		Version:       "1",
		Domain:        "acme.com",
		Namespace:     "namespace",
		DefaultLocale: language.English,
		Errors:        zedErrs,
	}

	err = exporter.Export(cfg, spec)
	require.NoError(t, err)

	buf, err := io.ReadAll(tempFile)
	require.NoError(t, err)

	expectedBuf := bytes.NewBuffer(nil)
	expectedBuf.Write(expectedGoErrorsFile)
	expectedBuf.Write(expectedGoEmbedFile)
	expectedBuf.Write(expectedLocaleEn)
	expectedBuf.Write(expectedLocaleRu)
	assert.Equal(t, expectedBuf.String(), string(buf))
}
