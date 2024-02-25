package input_test

import (
	"os"
	"testing"

	"github.com/amanbolat/zederr/internal/codegen/input"
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestYAMLUnmarshal(t *testing.T) {
	b, err := os.ReadFile("./testdata/01_valid_full_example.yaml")
	require.NoError(t, err)

	var spec input.ErrorListSpecification
	err = yaml.Unmarshal(b, &spec)
	require.NoError(t, err)

	pp.Println(spec)
	//
	// expected := []input.ErrorEntry{
	// 	{
	// 		ID:          "common.file_too_large",
	// 		GRPCCode:    3,
	// 		HTTPCode:    400,
	// 		description: "If the file received by server is larger than applied limit this error should be returned",
	// 		Translations: map[string]string{
	// 			"en":    "File is too large, max size is {{ .MaxSize | int }}",
	// 			"zh_cn": "上传的文件不能大于{{ .MaxSize | int }}",
	// 		},
	// 	},
	// 	{
	// 		ID:          "auth.unauthorized",
	// 		GRPCCode:    1,
	// 		HTTPCode:    403,
	// 		description: "User is not authorized to perform this action",
	// 		Translations: map[string]string{
	// 			"en":    "Please login to perform this action",
	// 			"zh_cn": "请登录再进行操作",
	// 		},
	// 	},
	// }

	// input.AssertErrorEntriesEquality(t, expected, entries)
}
