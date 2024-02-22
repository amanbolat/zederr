package core

import (
	"google.golang.org/grpc/codes"
)

type Error struct {
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

func (e Error) Code() string {
	return e.code
}

func (e Error) GrpcCode() codes.Code {
	return e.grpcCode
}

func (e Error) HttpCode() int {
	return e.httpCode
}

func (e Error) Description() string {
	return e.description
}

func (e Error) Translations() Localization {
	return e.localization
}

func (e Error) Arguments() []Argument {
	arr := make([]Argument, len(e.arguments))
	copy(arr, e.arguments)

	return arr
}
