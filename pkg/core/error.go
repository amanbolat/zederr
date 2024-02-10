package core

import (
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
)

type Language language.Tag

type Arguments map[string]interface{}

func (a Arguments) ToProtoStruct() (*structpb.Struct, error) {
	return structpb.NewStruct(a)
}

type LocalizedMessage struct {
	Lang Language
	Text string
}

type Error interface {
	error
	// UID returns a unique identifier of the error.
	// Format: <domain>/<namespace>/<name>
	UID() string
	Domain() string
	Namespace() string
	Code() string
	HTTPCode() int
	GRPCCode() codes.Code
	InternalMsg() string
	PublicMsg() string
	Args() Arguments
	// WithCause returns a new error with the given cause.
	// Usually it's used to wrap an application level error.
	WithCause(err error) Error
	Causes() []Error
}
