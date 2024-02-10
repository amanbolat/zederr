package goerr

import (
	"errors"
	"path"
	"reflect"
	"strings"

	"github.com/amanbolat/zederr/pkg/core"
	"google.golang.org/grpc/codes"
)

type Error struct {
	domain    string
	namespace string
	code      string
	err       error
	causes    []core.Error
}

func NewError(err error) core.Error {
	if err == nil {
		return nil
	}

	typ := reflect.TypeOf(err)
	pkgPath := getPkgPath(typ)

	s := typ.String()
	i := len(s) - 1
	sqBrackets := 0
	for i >= 0 && (s[i] != '.' || sqBrackets != 0) {
		switch s[i] {
		case ']':
			sqBrackets++
		case '[':
			sqBrackets--
		}
		i--
	}
	code := s[i+1:]

	var domain string
	var namespace string

	arr := strings.SplitN(pkgPath, "/", 2)
	if len(arr) == 2 {
		domain = arr[0]
		namespace = arr[1]
	}

	if len(arr) == 1 {
		domain = "golang.org"
		namespace = arr[0]
	}

	e := &Error{
		domain:    domain,
		namespace: namespace,
		code:      code,
		err:       err,
	}

	return e
}

func (e *Error) InternalMsg() string {
	return e.Error()
}

func (e *Error) PublicMsg() string {
	return ""
}

func (e *Error) Error() string {
	b := strings.Builder{}
	b.WriteString(e.UID())
	b.WriteString(": ")
	b.WriteString(e.err.Error())

	if len(e.causes) > 0 {
		b.WriteString("\ncauses:")
		for _, cause := range e.causes {
			b.WriteString(cause.Error())
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (e Error) UID() string {
	return path.Join(e.domain, e.namespace, e.code)
}

func (e Error) Domain() string {
	return e.domain
}

func (e Error) Namespace() string {
	return e.namespace
}

func (e Error) Code() string {
	return e.code
}

func (e Error) Description() string {
	return ""
}

func (e Error) HTTPCode() int {
	return 500
}

func (e Error) GRPCCode() codes.Code {
	return codes.Internal
}

func (e Error) Args() core.Arguments {
	return nil
}

func (e Error) Causes() []core.Error {
	return e.causes
}

func (e *Error) WithCause(err error) core.Error {
	if err == nil {
		return e
	}

	var v core.Error
	if errors.As(err, &v) {
		v = err.(core.Error)
	} else {
		v = NewError(err)
	}

	e.causes = append(e.causes, v)

	return e
}

func getPkgPath(t reflect.Type) string {
	pkgPath := t.PkgPath()
	if pkgPath != "" {
		return pkgPath
	}

	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		return getPkgPath(t.Elem())
	default:
		return ""
	}
}
