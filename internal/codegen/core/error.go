package core

import (
	"google.golang.org/grpc/codes"
)

type Error struct {
	id           string
	grpcCode     codes.Code
	httpCode     int
	description  string
	message      string
	isDeprecated bool
	localization Localization
	arguments    []Argument
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

func (e Error) Description() string {
	return e.description
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Localization() Localization {
	return e.localization
}

func (e Error) IsDeprecated() bool {
	return e.isDeprecated
}

func (e Error) Arguments() []Argument {
	arr := make([]Argument, len(e.arguments))
	copy(arr, e.arguments)

	return arr
}
