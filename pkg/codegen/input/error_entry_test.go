package input_test

import (
	"os"
	"testing"

	"github.com/amanbolat/zederr/codegen/input"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestYAMLUnmarshal(t *testing.T) {
	b, err := os.ReadFile("./testdata/01_two_valid_error_entries.yaml")
	require.NoError(t, err)

	var entries input.ErrorEntries
	err = yaml.Unmarshal(b, &entries)
	require.NoError(t, err)

	expected := []input.ErrorEntry{
		{
			ID:          "common.file_too_large",
			GRPCCode:    3,
			HTTPCode:    400,
			Description: "If the file received by server is larger than applied limit this error should be returned",
			Translations: map[string]string{
				"en":    "File is too large, max size is {{ .MaxSize | int }}",
				"zh_cn": "上传的文件不能大于{{ .MaxSize | int }}",
			},
		},
		{
			ID:          "auth.unauthorized",
			GRPCCode:    1,
			HTTPCode:    403,
			Description: "User is not authorized to perform this action",
			Translations: map[string]string{
				"en":    "Please login to perform this action",
				"zh_cn": "请登录再进行操作",
			},
		},
	}

	input.AssertErrorEntriesEquality(t, expected, entries)
}
