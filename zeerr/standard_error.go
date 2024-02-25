package zeerr

import (
	"bytes"
	"context"
	"html/template"
	"path"

	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
)

// standardError represents standardized error.
type standardError struct {
	code      string
	domain    string
	namespace string
	uid       string

	httpCode int
	grpcCode codes.Code
	args     Arguments

	internalMsg string
	publicMsg   string

	causes []Error
}

func NewError(
	ctx context.Context,
	localizer Localizer,
	code string,
	domain string,
	namespace string,
	httpCode int,
	grpcCode codes.Code,
	internalMsgTmpl *template.Template,
	args Arguments,
) Error {
	lang, ok := ctx.Value(LocaleCtxKeyType{}).(language.Tag)
	if !ok {
		lang = language.Und
	}

	uid := makeErrorUID(domain, namespace, code)

	publicMsg := localizer.LocalizePublicMessage(uid, lang, args)

	internalMsgBuf := bytes.NewBuffer(nil)
	_ = internalMsgTmpl.Execute(internalMsgBuf, args)

	return &standardError{
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

func (e standardError) UID() string {
	return e.uid
}

func (e standardError) Domain() string {
	return e.domain
}

func (e standardError) Namespace() string {
	return e.namespace
}

func (e standardError) Code() string {
	return e.code
}

// GRPCCode returns a gRPC code.
func (e standardError) GRPCCode() codes.Code {
	return e.grpcCode
}

// HTTPCode returns an HTTP code.
func (e standardError) HTTPCode() int {
	return e.httpCode
}

func (e standardError) InternalMsg() string {
	return e.internalMsg
}

func (e standardError) PublicMsg() string {
	return e.publicMsg
}

func (e standardError) Args() Arguments {
	return e.args
}

func (e standardError) Causes() []Error {
	return e.causes
}

func (e *standardError) WithCauses(causes ...Error) Error {
	for _, c := range causes {
		if c != nil {
			e.causes = append(e.causes, c)
		}
	}

	return e
}

// standardError implements error interface.
// NOTE: it returns the internal message and all the causes.
// It's not recommended to return the value to the client.
func (e *standardError) Error() string {
	return e.formattedErr()
}

func (e *standardError) formattedErr() string {
	buf := bytes.NewBuffer([]byte(e.internalMsg))

	for _, cause := range e.causes {
		buf.WriteString("\n\t")
		buf.WriteString(cause.Error())
	}

	return buf.String()
}
