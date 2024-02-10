package core

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/utf8string"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
)

type ErrorBuilder struct {
	p Parser
}

func NewErrorBuilder(p Parser) *ErrorBuilder {
	return &ErrorBuilder{
		p: p,
	}
}

func (b *ErrorBuilder) NewError(
	id string,
	grpcCode codes.Code,
	httpCode int,
	description string,
	translations map[string]string,
) (Error, error) {
	id = strings.TrimSpace(id)
	description = strings.TrimSpace(description)

	if id == "" {
		return Error{}, fmt.Errorf("error id is empty")
	}

	idArr := strings.Split(id, ".")
	if len(idArr) < 2 {
		return Error{}, fmt.Errorf("error id format is wrong; got %s", id)
	}

	if grpcCode == codes.OK {
		return Error{}, fmt.Errorf("grpc code should not be OK; got %s for error with id %s", grpcCode.String(), id)
	}

	if grpcCode > codes.Unauthenticated {
		slog.Warn("grpc code is not in the range of standard grpc codes", slog.Uint64("grpc_code", uint64(grpcCode)))
	}

	if httpCode < 100 || httpCode > 599 {
		slog.Warn("http code is not in the range of standard http codes", slog.Uint64("http_code", uint64(httpCode)))
	}

	if description == "" {
		slog.Warn("description is empty", slog.String("id", id))
	}

	if translations == nil {
		return Error{}, fmt.Errorf("translation are not provided; error id %s", id)
	}

	newTranslationsMap := make(map[string]string)

	fields := map[string]Param{}

	for lang, txt := range translations {
		_, err := language.Parse(lang)
		if err != nil {
			return Error{}, fmt.Errorf("failed to parse language tag %s for error with id %s: %w", lang, id, err)
		}

		txt = strings.TrimSpace(txt)
		if txt == "" {
			return Error{}, fmt.Errorf("translation for language %s is empty for error with id %s", lang, id)
		}

		fieldMap, localizableMessage, err := b.p.Parse(txt)
		if err != nil {
			return Error{}, fmt.Errorf("failed to parse translation for language %s for error with id %s: %w", lang, id, err)
		}

		for _, v := range fieldMap {
			if i, ok := fields[v.Name]; ok && i.Type != v.Type {
				return Error{}, fmt.Errorf("field %s of error with id %s has different types in multiple translations", v.Name, id)
			}

			v.Name = strings.TrimSpace(v.Name)
			if v.Name == "" {
				return Error{}, fmt.Errorf("found empty field name for error with id %s", id)
			}

			if !utf8string.NewString(v.Name).IsASCII() {
				return Error{}, fmt.Errorf("field name %s of error with id %s is not ASCII string", v.Name, id)
			}

			if !unicode.IsLetter([]rune(v.Name)[0]) {
				return Error{}, fmt.Errorf("field name %s of error with id %s should start with letter", v.Name, id)
			}

			v.Name = strcase.ToCamel(v.Name)

			fields[v.Name] = v
		}

		newTranslationsMap[lang] = localizableMessage
	}

	fieldsArr := maps.Values(fields)
	slices.SortFunc(fieldsArr, func(i, j Param) int {
		return strings.Compare(i.Name, j.Name)
	})

	return Error{
		id:           id,
		grpcCode:     grpcCode,
		httpCode:     httpCode,
		description:  description,
		translations: newTranslationsMap,
		fields:       fieldsArr,
	}, nil
}
