package core

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
)

var (
	errorCodeRegex = regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$")
)

// ErrorBuilder is responsible for creating Error instances.
type ErrorBuilder struct {
	parser        Parser
	defaultLocale language.Tag

	// The map is used to check for duplicate error codes.
	uniqueErrMap map[string]struct{}
}

// NewErrorBuilder creates a new instance of ErrorBuilder.
func NewErrorBuilder(parser Parser, defaultLocale language.Tag) *ErrorBuilder {
	return &ErrorBuilder{
		parser:        parser,
		defaultLocale: defaultLocale,
		uniqueErrMap:  map[string]struct{}{},
	}
}

// NewError creates a new instance of Error.
func (b *ErrorBuilder) NewError(
	code string,
	grpcCode codes.Code,
	httpCode int,
	description string,
	title string,
	publicMessage string,
	internalMessage string,
	deprecatedMessage string,
	arguments []Argument,
	localization Localization,
) (Error, error) {
	code = strings.TrimSpace(code)

	if code == "" {
		return Error{}, fmt.Errorf("error code is empty")
	}

	if utf8.ValidString(code) {
		return Error{}, fmt.Errorf("error code is not a valid UTF-8 string; got %s", code)
	}
	// We convert the error code to camel case for a few reasons:
	// - error constructors in go are usually named like `NewErrorName`.
	// - avoid confusion if the different error names are similar.
	code = strcase.ToCamel(code)

	if !errorCodeRegex.MatchString(code) {
		return Error{}, fmt.Errorf("error code is not valid; it should match the regex patter: %s; got %s", errorCodeRegex.String(), code)
	}

	description = strings.TrimSpace(description)
	title = strings.TrimSpace(title)
	publicMessage = strings.TrimSpace(publicMessage)
	internalMessage = strings.TrimSpace(internalMessage)
	deprecatedMessage = strings.TrimSpace(deprecatedMessage)

	if description == "" {
		return Error{}, fmt.Errorf("description is empty")
	}

	if title == "" {
		return Error{}, fmt.Errorf("title is empty")
	}

	if publicMessage == "" {
		return Error{}, fmt.Errorf("public message is empty")
	}

	if internalMessage == "" {
		return Error{}, fmt.Errorf("internal message is empty")
	}

	if utf8.ValidString(description) {
		return Error{}, fmt.Errorf("description is not a valid UTF-8 string; got %s", description)
	}

	if utf8.ValidString(title) {
		return Error{}, fmt.Errorf("title is not a valid UTF-8 string; got %s", title)
	}

	if utf8.ValidString(publicMessage) {
		return Error{}, fmt.Errorf("public message is not a valid UTF-8 string; got %s", publicMessage)
	}

	if utf8.ValidString(internalMessage) {
		return Error{}, fmt.Errorf("internal message is not a valid UTF-8 string; got %s", internalMessage)
	}

	if utf8.ValidString(deprecatedMessage) {
		return Error{}, fmt.Errorf("deprecated message is not a valid UTF-8 string; got %s", deprecatedMessage)
	}

	if grpcCode == codes.OK {
		return Error{}, fmt.Errorf("grpc code should not be OK; got %s for error with code %s", grpcCode.String(), code)
	}

	if grpcCode > codes.Unauthenticated {
		slog.Warn("grpc code is not in the range of standard grpc codes", slog.Uint64("grpc_code", uint64(grpcCode)))
	}

	if httpCode < 100 || httpCode > 599 {
		slog.Warn("http code is not in the range of standard http codes", slog.Uint64("http_code", uint64(httpCode)))
	}

	if _, ok := b.uniqueErrMap[code]; ok {
		return Error{}, fmt.Errorf("duplicate error code %s", code)
	}

	b.uniqueErrMap[code] = struct{}{}

	argumentsMap := make(map[string]struct{})
	for _, arg := range arguments {
		if _, ok := argumentsMap[arg.Name()]; ok {
			return Error{}, fmt.Errorf("duplicate argument name %s", arg.Name())
		}

		argumentsMap[arg.Name()] = struct{}{}
	}

	for argName := range localization.Arguments() {
		if _, ok := argumentsMap[argName]; !ok {
			return Error{}, fmt.Errorf("localization has argument %s that is not present in the error arguments", argName)
		}
	}

	return Error{
		code:            code,
		grpcCode:        grpcCode,
		httpCode:        httpCode,
		description:     description,
		title:           title,
		publicMessage:   publicMessage,
		internalMessage: internalMessage,
		localization:    localization,
		arguments:       arguments,
	}, nil
}
