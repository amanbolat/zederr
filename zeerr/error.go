package zeerr

import (
	"google.golang.org/grpc/codes"
)

type Arguments map[string]interface{}

type Error interface {
	UID() string
	Domain() string
	Namespace() string
	Code() string
	GRPCCode() codes.Code
	HTTPCode() int
	InternalMsg() string
	PublicMsg() string
	Args() Arguments
	Causes() []Error
	WithCauses(causes ...Error) Error
	Error() string
}
