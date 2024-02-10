package goerr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/amanbolat/zederr/pkg/goerr"
)

type testError struct {
	Msg string
}

var testErr = errors.New("test error")

func (e *testError) Error() string {
	return e.Msg
}

func TestNewError(t *testing.T) {
	err := goerr.NewError(&testError{Msg: "failure"})

	err = err.WithCause(testErr)
	// fmt.Println(err.Error())

	err2 := goerr.NewError(errors.New("cause 2"))
	err2 = err2.WithCause(errors.New("cause 3"))

	err = err.WithCause(err2)
	fmt.Println(err.Error())
	// fmt.Println(err2.Error())
}
