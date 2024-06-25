package zeerr

import (
	"bytes"
	"context"

	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
)

// Error represents a standardized error.
type Error struct {
	id          string
	httpCode    int
	grpcCode    codes.Code
	arguments   map[string]any
	message     string
	causes      []*Error
	internalErr error
}

func RestoreError(
	id string,
	httpCode int,
	grpcCode codes.Code,
	arguments map[string]any,
	message string,
	causes []*Error,
) *Error {
	return &Error{
		id:        id,
		httpCode:  httpCode,
		grpcCode:  grpcCode,
		arguments: arguments,
		message:   message,
		causes:    causes,
	}
}

// NewError creates a new Error.
func NewError(
	ctx context.Context,
	localizer Localizer,
	id string,
	httpCode int,
	grpcCode codes.Code,
	arguments map[string]any,
) *Error {
	lang, ok := ctx.Value(LocaleCtxKeyType{}).(language.Tag)
	if !ok {
		lang = language.Und
	}

	publicMsg := localizer.LocalizeMessage(id, lang, arguments)

	return &Error{
		id:          id,
		httpCode:    httpCode,
		grpcCode:    grpcCode,
		arguments:   arguments,
		message:     publicMsg,
		causes:      nil,
		internalErr: nil,
	}
}

func (e Error) ID() string {
	return e.id
}

func (e Error) GRPCCode() codes.Code {
	return e.grpcCode
}

func (e Error) HTTPCode() int {
	return e.httpCode
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Arguments() map[string]any {
	return e.arguments
}

func (e Error) Causes() []*Error {
	return e.causes
}

func (e Error) InternalErr() error {
	return e.internalErr
}

func (e *Error) WithCauses(causes ...*Error) *Error {
	for _, c := range causes {
		if c != nil {
			e.causes = append(e.causes, c)
		}
	}

	return e
}

func (e *Error) WithInternalError(err error) *Error {
	e.internalErr = err

	return e
}

func (e *Error) Error() string {
	return e.formattedErr()
}

func (e *Error) formattedErr() string {
	buf := bytes.NewBuffer([]byte(e.message))

	for _, cause := range e.causes {
		buf.WriteString("\n\t")
		buf.WriteString(cause.Error())
	}

	return buf.String()
}
