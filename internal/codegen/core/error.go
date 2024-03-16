package core

import (
	"path"

	"google.golang.org/grpc/codes"
)

type Error struct {
	domain          string
	namespace       string
	code            string
	grpcCode        codes.Code
	httpCode        int
	description     string
	title           string
	publicMessage   string
	internalMessage string
	localization    Localization
	arguments       []Argument
}

func (e Error) UID() string {
	return path.Join(e.domain, e.namespace, e.code)
}

func (e Error) Code() string {
	return e.code
}

func (e Error) Domain() string {
	return e.domain
}

func (e Error) Namespace() string {
	return e.namespace
}

func (e Error) GrpcCode() codes.Code {
	return e.grpcCode
}

func (e Error) HTTPCode() int {
	return e.httpCode
}

func (e Error) Description() string {
	return e.description
}

func (e Error) Title() string {
	return e.title
}

func (e Error) PublicMessage() string {
	return e.publicMessage
}

func (e Error) InternalMessage() string {
	return e.internalMessage
}

func (e Error) Translations() Localization {
	return e.localization
}

func (e Error) Arguments() []Argument {
	arr := make([]Argument, len(e.arguments))
	copy(arr, e.arguments)

	return arr
}
