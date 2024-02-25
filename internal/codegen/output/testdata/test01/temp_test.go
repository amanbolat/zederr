package zederrtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/amanbolat/zederr/zeerr"
)

func TestName(t *testing.T) {
	err := NewCode(context.Background(), "ARGUMENT 1")
	fmt.Println(err.(zeerr.Error).PublicMsg())
}

var err error

func BenchmarkName(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err = NewCode(context.Background(), "error")
	}
}
