package core

import (
	"bytes"
	"context"
	"html/template"
	"path"

	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
)

type LocaleCtxKeyType struct{}

type Arguments map[string]interface{}

// Error represents standardized error.
type Error struct {
	code      string
	domain    string
	namespace string
	uid       string

	httpCode int
	grpcCode codes.Code
	args     Arguments

	internalMsg string
	publicMsg   string

	causes []*Error
}

func NewLocalizedError(
	ctx context.Context,
	localizer Localizer,
	code string,
	domain string,
	namespace string,
	httpCode int,
	grpcCode codes.Code,
	internalMsgTmpl *template.Template,
	args Arguments,
) *Error {
	lang, ok := ctx.Value(LocaleCtxKeyType{}).(language.Tag)
	if !ok {
		lang = language.Und
	}

	uid := makeErrorUID(domain, namespace, code)

	publicMsg := localizer.LocalizePublicMessage(uid, lang, args)

	// tmpl, _ := template.New("").Parse(internalMsgTmpl)
	internalMsgBuf := bytes.NewBuffer(nil)
	_ = internalMsgTmpl.Execute(internalMsgBuf, args)

	return &Error{
		code:        code,
		domain:      domain,
		namespace:   namespace,
		uid:         uid,
		httpCode:    httpCode,
		grpcCode:    grpcCode,
		args:        args,
		internalMsg: internalMsgBuf.String(),
		publicMsg:   publicMsg,
		causes:      nil,
	}
}

func makeErrorUID(domain, namespace, code string) string {
	return path.Join(domain, namespace, code)
}

func NewError(
	code string,
	domain string,
	namespace string,
	httpCode int,
	grpcCode codes.Code,
	publicMsg string,
	internalMsg string,
	args Arguments,
) *Error {
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

func (e Error) UID() string {
	return e.uid
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

func (e *Error) Args() Arguments {
	return e.args
}

func (e *Error) Causes() []*Error {
	return e.causes
}

func (e *Error) WithCauses(causes ...*Error) *Error {
	for _, c := range causes {
		if c != nil {
			e.causes = append(e.causes, c)
		}
	}

	return e
}

// Error implements error interface.
// NOTE: it returns the internal message and all the causes.
// It's not recommended to return the value to the client.
func (e *Error) Error() string {
	return e.formattedErr()
}

func (e *Error) formattedErr() string {
	buf := bytes.NewBuffer([]byte(e.internalMsg))

	for _, cause := range e.causes {
		buf.WriteString("\n\t")
		buf.WriteString(cause.Error())
	}

	return buf.String()
}
