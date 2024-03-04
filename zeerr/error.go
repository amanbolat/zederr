package zeerr

import (
	"google.golang.org/grpc/codes"
)

// Arguments represents Error arguments.
type Arguments map[string]interface{}

// Error represents standardized error.
//
// The interface should not be implemented directly.
// Use generated errors and constructors to create new errors instances.
type Error interface {
	// UID is a unique identifier of the error.
	UID() string
	// Domain represents error's domain.
	// Usually it's a domain owned by the organization.
	Domain() string
	// Namespace can represent a service, package, or module.
	// For example all the errors under the "auth" namespaces
	// are related to the authentication and authorization.
	Namespace() string
	// Code returns error's code.
	Code() string
	// GRPCCode is a gRPC status code that will be returned to the
	// client in case the error was returned via gRPC.
	GRPCCode() codes.Code
	// HTTPCode is an HTTP status code that will be returned to the
	// client in case the error was returned via HTTP.
	HTTPCode() int
	// InternalMsg is used to describe the error and not exposed to
	// external consumers of the service.
	InternalMsg() string
	// PublicMsg provides a user-friendly message, which is localized based
	// on the user's locale and will be returned to the client.
	PublicMsg() string
	// Args returns error's arguments.
	Args() Arguments
	// Causes returns a list of errors that caused the current error.
	// It's usually used to nest other errors.
	// Image the usecase, when the user submits a complex form and the latter
	// contains multiple sections with different fields.
	// Sometimes, it is required to validate each field and return errors
	// in a structured manner, so they can be easily rendered on the client side.
	Causes() []Error
	// WithCauses is used to attach causes to the error.
	WithCauses(causes ...Error) Error
	// Error implements the error interface.
	Error() string
}
