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

var errorCodeRegex = regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$")

// ErrorBuilder is responsible for creating Error instances.
type ErrorBuilder struct {
	defaultLocale language.Tag

	// The map is used to check for duplicate error codes.
	uniqueErrMap map[string]struct{}
}

// NewErrorBuilder creates a new instance of ErrorBuilder.
func NewErrorBuilder(specVersion, defaultLocale string) (*ErrorBuilder, error) {
	locale, err := language.Parse(defaultLocale)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default locale; %w", err)
	}

	if specVersion != "1" {
		return nil, fmt.Errorf("spec version is not supported; got %s", specVersion)
	}

	return &ErrorBuilder{
		defaultLocale: locale,
		uniqueErrMap:  map[string]struct{}{},
	}, nil
}

// NewError creates a new instance of Error.
func (b *ErrorBuilder) NewError(
	id string,
	message string,
	grpcCode codes.Code,
	httpCode int,
	description string,
	isDeprecated bool,
	arguments []Argument,
	localization Localization,
) (Error, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return Error{}, fmt.Errorf("error code is empty")
	}

	if !utf8.ValidString(id) {
		return Error{}, fmt.Errorf("error code is not a valid UTF-8 string; got %s", id)
	}
	// We convert the error id to camel case for a few reasons:
	// - error constructors in go are usually named like `NewErrorName`.
	// - avoid confusion if the different error names are similar.
	id = strcase.ToSnake(id)

	if !errorCodeRegex.MatchString(id) {
		return Error{}, fmt.Errorf("error id is not valid; it should match the regex patter: %s; got %s", errorCodeRegex.String(), id)
	}

	description = strings.TrimSpace(description)
	message = strings.TrimSpace(message)

	if message == "" {
		return Error{}, fmt.Errorf("public message is empty")
	}

	if !utf8.ValidString(message) {
		return Error{}, fmt.Errorf("public message is not a valid UTF-8 string; got %s", message)
	}

	if description == "" {
		return Error{}, fmt.Errorf("description is empty")
	}

	if !utf8.ValidString(description) {
		return Error{}, fmt.Errorf("description is not a valid UTF-8 string; got %s", description)
	}

	if grpcCode == codes.OK {
		return Error{}, fmt.Errorf("grpc code should not be OK; got %s for error with code %s", grpcCode.String(), id)
	}

	if grpcCode > codes.Unauthenticated {
		slog.Warn("grpc code is not in the range of standard grpc codes", slog.Uint64("grpc_code", uint64(grpcCode)))
	}

	if httpCode < 100 || httpCode > 599 {
		slog.Warn("http code is not in the range of standard http codes", slog.Uint64("http_code", uint64(httpCode)))
	}

	if _, ok := b.uniqueErrMap[id]; ok {
		return Error{}, fmt.Errorf("duplicate error code %s", id)
	}

	b.uniqueErrMap[id] = struct{}{}

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

	templateValidator := NewTemplateValidator(&TemplateValidatorConfig{
		Debug:     false,
		Arguments: argumentsMap,
	})

	err := templateValidator.Validate(message)
	if err != nil {
		return Error{}, fmt.Errorf("public message is not a valid template; %w", err)
	}

	for lang, msg := range localization.Message() {
		err := templateValidator.Validate(msg)
		if err != nil {
			return Error{}, fmt.Errorf("public message for %s language is not a valid template; %w", lang, err)
		}
	}

	for argName, translations := range localization.Arguments() {
		for lang, msg := range translations {
			err := templateValidator.Validate(msg)
			if err != nil {
				return Error{}, fmt.Errorf("argument %s has an invalid localized message template for %s language; %w", argName, lang, err)
			}
		}
	}

	// All the parameters passed to the constructor are considered as text in default locale,
	// therefore `localization` should be propagated with them.
	err = localization.AddMessageTranslation(b.defaultLocale.String(), message)
	if err != nil {
		return Error{}, fmt.Errorf("failed to add public message translation for default locale; %w", err)
	}

	err = localization.AddDescriptionTranslation(b.defaultLocale.String(), description)
	if err != nil {
		return Error{}, fmt.Errorf("failed to add description translation for default locale; %w", err)
	}

	for _, arg := range arguments {
		err = localization.AddArgumentTranslation(arg.Name(), arg.Description(), b.defaultLocale.String())
		if err != nil {
			return Error{}, fmt.Errorf("failed to add argument (%s) translation for default locale; %w", arg.Name(), err)
		}
	}

	return Error{
		id:           id,
		grpcCode:     grpcCode,
		httpCode:     httpCode,
		description:  description,
		message:      message,
		isDeprecated: isDeprecated,
		localization: localization,
		arguments:    arguments,
	}, nil
}
