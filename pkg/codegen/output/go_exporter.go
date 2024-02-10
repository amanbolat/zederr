package output

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "embed"

	"github.com/BurntSushi/toml"
	"github.com/amanbolat/zederr/pkg/codegen/core"
	"github.com/iancoleman/strcase"
)

var (
	// keywords is a map of Go keywords, package names, builtins and param names that should be renamed.
	keywords = map[string]struct{}{
		// Go keywords.
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

//go:embed templates/go_embed.tmpl
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

// localeEntry represents a single translation entry in a locale file.
//
// Example toml representation:
//
//	["error.auth.unauthorized"] <-- provided by map key
//	other = "请登录再进行操作" <-- localeEntry
type localeEntry struct {
	Other string `toml:"other"`
}

type GoExporter struct{}

func NewGoExporter() *GoExporter {
	return &GoExporter{}
}

func (e *GoExporter) Export(cfg core.GoExporterConfig, errors []core.Error) error {
	if cfg.OutputPath != "" {
		err := os.MkdirAll(cfg.OutputPath, 0755)
		if err != nil {
			return fmt.Errorf("GoExporter: failed to create output directory: %w", err)
		}
	}

	err := e.renderErrors(cfg, errors)
	if err != nil {
		return fmt.Errorf("failed to render errors: %w", err)
	}

	err = e.renderLocalesEmbed(cfg, errors)
	if err != nil {
		return fmt.Errorf("failed to render locales embed: %w", err)
	}

	err = e.renderLocales(cfg, errors)
	if err != nil {
		return fmt.Errorf("failed to render locales: %w", err)
	}

	return nil
}

func (e *GoExporter) renderLocalesEmbed(cfg core.GoExporterConfig, errors []core.Error) error {
	locales := make(map[string]struct{})

	for _, coreErr := range errors {
		for locale := range coreErr.Translations() {
			locales[locale] = struct{}{}
		}
	}

	var localesTemplateData []localeTemplateData
	for locale := range locales {
		localesTemplateData = append(localesTemplateData, localeTemplateData{
			FileName: fmt.Sprintf("locale.%s.toml", locale),
			Lang:     locale,
		})
	}

	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"toUpper": strings.ToUpper,
	})

	_, err := tmpl.Parse(goErrorLocalesEmbed)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, goErrorLocalesTemplateData{
		PackageName: cfg.PackageName,
		Locales:     localesTemplateData,
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated go code: %w", err)
	}

	outputBuf := bytes.NewReader(formattedSource)

	if cfg.Output != nil {
		_, err = io.Copy(cfg.Output, outputBuf)
		if err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	}

	if cfg.OutputPath != "" {
		fileName := filepath.Join(cfg.OutputPath, "error_locales_embed.go")

		err = os.WriteFile(fileName, formattedSource, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *GoExporter) renderErrors(cfg core.GoExporterConfig, errors []core.Error) error {
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"errorConstructorParams": errorConstructorParams,
		"toLowerCamel":           strcase.ToLowerCamel,
		"toParamName":            toParamName,
		"toUpper":                strings.ToUpper,
	})

	_, err := tmpl.Parse(goErrorsTemplate)
	if err != nil {
		return fmt.Errorf("GoExporter: failed to parse template: %w", err)
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
		return fmt.Errorf("failed to execute template: %w", err)
	}

	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated go code: %w", err)
	}

	outputBuf := bytes.NewReader(formattedSource)

	if cfg.Output != nil {
		_, err = io.Copy(cfg.Output, outputBuf)
		if err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	}

	if cfg.OutputPath != "" {
		fileName := filepath.Join(cfg.OutputPath, "errors.go")

		err = os.WriteFile(fileName, formattedSource, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *GoExporter) renderLocales(cfg core.GoExporterConfig, errors []core.Error) error {
	errMaps := map[string]map[string]localeEntry{}
	for _, coreErr := range errors {
		for lang, translation := range coreErr.Translations() {
			if _, ok := errMaps[lang]; !ok {
				errMaps[lang] = make(map[string]localeEntry)
			}
			errMaps[lang][coreErr.ID()] = localeEntry{Other: translation}
		}
	}

	for lang, v := range errMaps {
		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		enc.Indent = ""
		err := enc.Encode(v)
		if err != nil {
			return fmt.Errorf("failed to encode toml: %w", err)
		}

		if cfg.Output != nil {
			_, err = io.Copy(cfg.Output, &buf)
			if err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}

		if cfg.OutputPath != "" {
			fileName := filepath.Join(cfg.OutputPath, fmt.Sprintf("locale.%s.toml", lang))

			err = os.WriteFile(fileName, buf.Bytes(), 0666)
			if err != nil {
				return fmt.Errorf("failed to write %s error locale messages to file: %w", lang, err)
			}
		}
	}

	return nil
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
