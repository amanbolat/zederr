package core

import (
	"strings"

	"github.com/iancoleman/strcase"
	"google.golang.org/grpc/codes"
)

type Error struct {
	id           string
	grpcCode     codes.Code
	httpCode     int
	description  string
	translations map[string]string
	fields       []Param
}

func (e Error) ID() string {
	return e.id
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

func (e Error) Translations() map[string]string {
	m := make(map[string]string, len(e.translations))
	for k, v := range e.translations {
		m[k] = v
	}

	return m
}

func (e Error) Name() string {
	arr := strings.Split(e.id, ".")
	if len(arr) == 0 {
		return ""
	}

	return strcase.ToCamel(strings.Join(arr, "_"))
}

func (e Error) Fields() []Param {
	arr := make([]Param, len(e.fields))
	copy(arr, e.fields)

	return arr
}
