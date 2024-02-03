package output

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"

	_ "embed"

	"github.com/amanbolat/zederr/internal/core"
	"github.com/iancoleman/strcase"
)

var (
	// keywords is a map of Go keywords, package names, builtins and param names that should be renamed.
	keywords = map[string]struct{}{
		"break":       {},
		"default":     {},
		"func":        {},
		"interface":   {},
		"select":      {},
		"case":        {},
		"defer":       {},
		"go":          {},
		"map":         {},
		"struct":      {},
		"chan":        {},
		"else":        {},
		"goto":        {},
		"package":     {},
		"switch":      {},
		"const":       {},
		"fallthrough": {},
		"if":          {},
		"range":       {},
		"type":        {},
		"continue":    {},
		"for":         {},
		"import":      {},
		"return":      {},
		"var":         {},
		"error":       {},
		// Parameter names.
		"err": {},
		// Package aliases.
		"pkgzederr":   {},
		"pkgcodes":    {},
		"pkgstructpb": {},
	}
)

//go:embed templates/go_errors.tmpl
var goErrorsTemplate string

//go:embed templates/go_error_locales_embed.tmpl
var goErrorLocalesEmbed string

type goErrorsTemplateData struct {
	PackageName string
	Imports     []string
	Errors      []core.Error
}

type localeTemplateData struct {
	FileName string
	Lang     string
}

type goErrorLocalesTemplateData struct {
	PackageName string
	Locales     []localeTemplateData
}

type GoExporter struct{}

func NewGoExporter() *GoExporter {
	return &GoExporter{}
}

func (e *GoExporter) Export(cfg core.GoExporterConfig, errors []core.Error) error {
	renderedErrors, err := e.renderErrors(cfg, errors)
	if err != nil {
		return fmt.Errorf("GoExporter: failed to render errors: %w", err)
	}

	renderedLocalesEmbed, err := e.renderLocalesEmbed(cfg, errors)
	if err != nil {
		return fmt.Errorf("GoExporter: failed to render locales embed: %w", err)
	}

	if cfg.Output != nil {
		_, err = io.Copy(cfg.Output, renderedErrors)
		if err != nil {
			return fmt.Errorf("GoExporter: failed to write to output: %w", err)
		}

		_, err = io.Copy(cfg.Output, renderedLocalesEmbed)
		if err != nil {
			return fmt.Errorf("GoExporter: failed to write to output: %w", err)
		}
	}

	if cfg.OutputPath != "" {
		err = os.MkdirAll(cfg.OutputPath, 0755)
		if err != nil {
			return fmt.Errorf("GoExporter: failed to create output directory: %w", err)
		}

		errorFile, err := os.Create(cfg.OutputPath + "/errors.go")
		if err != nil {
			return fmt.Errorf("GoExporter: failed to create output file: %w", err)
		}

		_, err = io.Copy(errorFile, renderedErrors)
		if err != nil {
			return fmt.Errorf("GoExporter: failed to write to output file: %w", err)
		}

		err = errorFile.Close()
		if err != nil {
			return fmt.Errorf("GoExporter: failed to close output file: %w", err)
		}

		localesEmbedFile, err := os.Create(cfg.OutputPath + "/error_locales_embed.go")
		if err != nil {
			return fmt.Errorf("GoExporter: failed to create output file: %w", err)
		}

		_, err = io.Copy(localesEmbedFile, renderedLocalesEmbed)
		if err != nil {
			return fmt.Errorf("GoExporter: failed to write to output file: %w", err)
		}

		err = localesEmbedFile.Close()
		if err != nil {
			return fmt.Errorf("GoExporter: failed to close output file: %w", err)
		}
	}

	return nil
}

func (e *GoExporter) renderLocalesEmbed(cfg core.GoExporterConfig, errors []core.Error) (io.Reader, error) {
	locales := make(map[string]struct{})

	for _, coreErr := range errors {
		for locale := range coreErr.Translations() {
			locales[locale] = struct{}{}
		}
	}

	var localesTemplateData []localeTemplateData
	for locale := range locales {
		localesTemplateData = append(localesTemplateData, localeTemplateData{
			FileName: fmt.Sprintf("error_locales.%s.toml", locale),
			Lang:     locale,
		})
	}

	tmpl := template.New("")

	_, err := tmpl.Parse(goErrorLocalesEmbed)
	if err != nil {
		return nil, fmt.Errorf("GoExporter: failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, goErrorLocalesTemplateData{
		PackageName: cfg.PackageName,
		Locales:     localesTemplateData,
	})
	if err != nil {
		return nil, fmt.Errorf("GoExporter: failed to execute template: %w", err)
	}

	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("GoExporter: failed to format source: %w", err)
	}

	outputBuf := bytes.NewReader(formattedSource)

	return outputBuf, nil
}

func (e *GoExporter) renderErrors(cfg core.GoExporterConfig, errors []core.Error) (io.Reader, error) {
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"errorConstructorParams": errorConstructorParams,
		"toLowerCamel":           strcase.ToLowerCamel,
		"toParamName":            toParamName,
	})

	_, err := tmpl.Parse(goErrorsTemplate)
	if err != nil {
		return nil, fmt.Errorf("GoExporter: failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, goErrorsTemplateData{
		PackageName: cfg.PackageName,
		Imports: []string{
			`pkgzederr "github.com/amanbolat/zederr/pkg/zederr"`,
			`pkgcodes "google.golang.org/grpc/codes"`,
			`pkgstructpb "google.golang.org/protobuf/types/known/structpb"`,
		},
		Errors: errors,
	})
	if err != nil {
		return nil, fmt.Errorf("GoExporter: failed to execute template: %w", err)
	}

	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("GoExporter: failed to format source: %w", err)
	}

	outputBuf := bytes.NewReader(formattedSource)

	return outputBuf, nil
}

func errorConstructorParams(e core.Error) string {
	var res string
	for i, field := range e.Fields() {
		name := toParamName(field.Name)
		res += name + " " + field.Type
		if i != len(e.Fields())-1 {
			res += ", "
		}
	}

	return res
}

func toParamName(name string) string {
	name = strcase.ToSnake(name)

	_, ok := keywords[name]
	if ok {
		return "param" + strcase.ToCamel(name)
	}

	return name
}
