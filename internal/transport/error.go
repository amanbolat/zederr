package transport

import (
	"github.com/amanbolat/zederr/zeerr"
	"google.golang.org/grpc/codes"
)

type transportError struct {
	code      string
	domain    string
	namespace string
	uid       string

	httpCode int
	grpcCode codes.Code
	args     zeerr.Arguments

	internalMsg string
	publicMsg   string

	causes []zeerr.Error
}

func NewError(
	code string,
	domain string,
	namespace string,
	uid string,
	httpCode int,
	grpcCode codes.Code,
	args zeerr.Arguments,
	internalMsg string,
	publicMsg string,
	causes []zeerr.Error,
) zeerr.Error {
	return &transportError{
		code:        code,
		domain:      domain,
		namespace:   namespace,
		uid:         uid,
		httpCode:    httpCode,
		grpcCode:    grpcCode,
		args:        args,
		internalMsg: internalMsg,
		publicMsg:   publicMsg,
		causes:      causes,
	}
}

func (e transportError) UID() string {
	return e.uid
}

func (e transportError) Domain() string {
	return e.domain
}

func (e transportError) Namespace() string {
	return e.namespace
}

func (e transportError) Code() string {
	return e.code
}

func (e transportError) GRPCCode() codes.Code {
	return e.grpcCode
}

func (e transportError) HTTPCode() int {
	return e.httpCode
}

func (e transportError) InternalMsg() string {
	return e.internalMsg
}

func (e transportError) PublicMsg() string {
	return e.publicMsg
}

func (e transportError) Args() zeerr.Arguments {
	return e.args
}

func (e transportError) Causes() []zeerr.Error {
	return e.causes
}

func (e *transportError) WithCauses(causes ...zeerr.Error) zeerr.Error {
	for _, c := range causes {
		if c != nil {
			e.causes = append(e.causes, c)
		}
	}

	return e
}

func (e transportError) Error() string {
	return e.internalMsg
}
