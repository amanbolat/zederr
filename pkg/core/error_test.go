package core_test

import (
	"fmt"
	"testing"

	"github.com/amanbolat/zederr/pkg/core"
)

func TestName(t *testing.T) {
	err1 := core.NewError(
		"code",
		"domain",
		"namespace",
		500,
		13,
		"PUBLIC {{ .attr_1 }}",
		"INTERNAL {{ .attr_1 }}",
		core.Arguments{
			"attr_1": 1,
		},
	)

	err2 := core.NewError(
		"code",
		"domain",
		"namespace",
		500,
		13,
		"PUBLIC {{ .a1 }}",
		"INTERNAL {{ .a1 }}",
		core.Arguments{
			"a1": 1,
		},
	)
	err1.WithCauses(err2)
	fmt.Println(err1.Error())
}
