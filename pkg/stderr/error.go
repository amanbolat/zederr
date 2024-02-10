package stderr

import (
	"bytes"
	"errors"
	"path"
	"text/template"

	"github.com/amanbolat/zederr/pkg/core"
	"github.com/amanbolat/zederr/pkg/goerr"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
)

type Language language.Tag

type LocalizedMessage struct {
	Lang Language
	Text string
}

// Error represents standardized error.
type Error struct {
	code      string
	domain    string
	namespace string

	httpCode int
	grpcCode codes.Code
	args     core.Arguments

	internalMsg string
	publicMsg   string

	causes []core.Error
}

func NewError(
	code string,
	domain string,
	namespace string,
	httpCode int,
	grpcCode codes.Code,
	internalMsg string,
	publicMsg string,
	args core.Arguments,
) core.Error {
	return &Error{
		code:        code,
		domain:      domain,
		namespace:   namespace,
		httpCode:    httpCode,
		grpcCode:    grpcCode,
		args:        args,
		internalMsg: internalMsg,
		publicMsg:   publicMsg,
		causes:      nil,
	}
}

func (e *Error) UID() string {
	return path.Join(e.domain, e.namespace, e.code)
}

func (e *Error) Domain() string {
	return e.domain
}

func (e *Error) Namespace() string {
	return e.namespace
}

func (e *Error) Code() string {
	return e.code
}

// GRPCCode returns a gRPC code.
func (e Error) GRPCCode() codes.Code {
	return e.grpcCode
}

// HTTPCode returns an HTTP code.
func (e Error) HTTPCode() int {
	return e.httpCode
}

func (e *Error) InternalMsg() string {
	return e.internalMsg
}

func (e *Error) PublicMsg() string {
	return e.publicMsg
}

func (e *Error) Args() core.Arguments {
	return e.args
}

func (e *Error) Causes() []core.Error {
	return e.causes
}

func (e *Error) WithCause(err error) core.Error {
	if err == nil {
		return e
	}

	var v core.Error
	if !errors.As(err, &v) {
		v = goerr.NewError(err)
	}

	e.causes = append(e.causes, v)

	return e
}

// Error implements error interface.
// NOTE: it returns the internal message and all the causes.
// It's not recommended to return the value to the client.
func (e *Error) Error() string {
	return e.formatted()
}

func (e *Error) formatted() string {
	tmpl, _ := template.New("").Parse(e.internalMsg)

	buf := bytes.NewBuffer(nil)
	_ = tmpl.Execute(buf, e.args)

	for _, cause := range e.causes {
		buf.WriteString("\n\t")
		buf.WriteString(cause.Error())
	}

	return buf.String()
}
