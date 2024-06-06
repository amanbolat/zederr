package zeerr

import (
	"bytes"
	"context"
	"html/template"

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

	internalErr error
}

// NewError creates a new error.
// NOTE: The constructor is meant to be used only by the generated code.
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

	uid := MakeUID(domain, namespace, code)

	publicMsg := localizer.LocalizePublicMessage(uid, lang, args)

	internalMsgBuf := bytes.NewBuffer(nil)
	err := internalMsgTmpl.Execute(internalMsgBuf, args)
	if err != nil {
		panic(err)
	}

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

// UID returns an error UID.
func (e standardError) UID() string {
	return e.uid
}

// Domain returns an error domain.
func (e standardError) Domain() string {
	return e.domain
}

// Namespace returns an error namespace.
func (e standardError) Namespace() string {
	return e.namespace
}

// Code returns an error code.
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

// InternalMsg returns an internal message.
func (e standardError) InternalMsg() string {
	return e.internalMsg
}

func (e standardError) PublicMsg() string {
	return e.publicMsg
}

// Args returns error's arguments.
func (e standardError) Args() Arguments {
	return e.args
}

// Causes returns a list of errors that caused the current error.
func (e standardError) Causes() []Error {
	return e.causes
}

// WithCauses is used to attach causes to the error.
// Causes should be of type Error and it's used to create
// complex nested error structures.
//
// Causes will not be unwrapped by the standard error implementation.
// To include internal error use WithInternalErr instead.
func (e *standardError) WithCauses(causes ...Error) Error {
	for _, c := range causes {
		if c != nil {
			e.causes = append(e.causes, c)
		}
	}

	return e
}

// WithInternalErr is used to attach an internal error to the error.
// This error will be unwrapped by the standard error implementation.
// Internal error will not be included in the encoded transport error.
func (e *standardError) WithInternalErr(internalErr error) Error {
	e.internalErr = internalErr

	return e
}

// Unwrap returns an internal error.
func (e *standardError) Unwrap() error {
	return e.internalErr
}

// Error implements the error interface.
func (e *standardError) Error() string {
	return e.formattedErr()
}

// formattedErr returns a formatted error message.
func (e *standardError) formattedErr() string {
	buf := bytes.NewBuffer([]byte(e.internalMsg))

	for _, cause := range e.causes {
		buf.WriteString("\n\t")
		buf.WriteString(cause.Error())
	}

	if e.internalErr != nil {
		buf.WriteString("\n\t")
		buf.WriteString("internal error: ")
		buf.WriteString(e.internalErr.Error())
	}

	return buf.String()
}
