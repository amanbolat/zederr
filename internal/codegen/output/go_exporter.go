package output

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/iancoleman/strcase"
	"golang.org/x/text/language"

	"github.com/amanbolat/zederr/internal/codegen/core"
)

// keywords is a map of Go keywords, package names, builtins and param names that should be renamed.
var keywords = map[string]struct{}{
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

//go:embed templates/go_errors.tmpl
var goErrorsTemplate string

//go:embed templates/go_embed.tmpl
var goErrorLocalesEmbed string

type goErrorsTemplateData struct {
	PackageName   string
	DefaultLocale string
	Imports       []string
	Errors        []core.Error
}

type localeTemplateData struct {
	FileName string
	Lang     language.Tag
}

type goErrorLocalesTemplateData struct {
	PackageName string
	Locales     []localeTemplateData
}

// localeEntry represents a single translation entry in a locale file.
//
// Example toml representation:
//
//	["acme.com/auth/unauthorized"] <-- provided by map key
//	other = "Please sign in" <-- localeEntry
type localeEntry struct {
	Other string `toml:"other"`
}

type GoExporter struct{}

func NewGoExporter() *GoExporter {
	return &GoExporter{}
}

func (e *GoExporter) Export(cfg core.ExportGo, spec core.Spec) error {
	if cfg.OutputPath == "" {
		return fmt.Errorf("output path is empty")
	}

	err := os.MkdirAll(cfg.OutputPath, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	err = e.renderErrors(cfg, spec)
	if err != nil {
		return fmt.Errorf("failed to render errors: %w", err)
	}

	err = e.renderLocalesEmbed(cfg, spec)
	if err != nil {
		return fmt.Errorf("failed to render locales embed: %w", err)
	}

	err = e.renderLocales(cfg, spec)
	if err != nil {
		return fmt.Errorf("failed to render locales: %w", err)
	}

	return nil
}

func (e *GoExporter) renderLocalesEmbed(cfg core.ExportGo, spec core.Spec) error {
	locales := make(map[language.Tag]struct{})

	for _, coreErr := range spec.Errors {
		for _, locale := range coreErr.Translations().AllLanguages() {
			locales[locale] = struct{}{}
		}
	}

	localesTemplateData := make([]localeTemplateData, 0, len(locales))
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

	fileName := filepath.Join(cfg.OutputPath, "error_locales_embed.go")

	err = os.WriteFile(fileName, formattedSource, 0o600)
	if err != nil {
		return err
	}

	return nil
}

func (e *GoExporter) renderErrors(cfg core.ExportGo, spec core.Spec) error {
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"errorConstructorParams": errorConstructorParams,
		"toLowerCamel":           strcase.ToLowerCamel,
		"toParamName":            toParamName,
		"toUpper":                strings.ToUpper,
	})

	_, err := tmpl.Parse(goErrorsTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, goErrorsTemplateData{
		PackageName:   cfg.PackageName,
		DefaultLocale: spec.DefaultLocale.String(),
		Imports: []string{
			`context "context"`,
			`template "html/template"`,
			`zeerr "github.com/amanbolat/zederr/zeerr"`,
			`zei18n "github.com/amanbolat/zederr/zei18n"`,
			`pkgcodes "google.golang.org/grpc/codes"`,
			`time "time"`,
		},
		Errors: spec.Errors,
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated go code: %w", err)
	}

	fileName := filepath.Join(cfg.OutputPath, "errors.go")

	err = os.WriteFile(fileName, formattedSource, 0o600)
	if err != nil {
		return err
	}

	return nil
}

func (e *GoExporter) renderLocales(cfg core.ExportGo, spec core.Spec) error {
	entryMap := map[language.Tag]map[string]localeEntry{}

	for _, coreErr := range spec.Errors {
		for _, lang := range coreErr.Translations().AllLanguages() {
			if _, ok := entryMap[lang]; !ok {
				entryMap[lang] = make(map[string]localeEntry)
			}
		}
	}

	for _, coreErr := range spec.Errors {
		for lang, translation := range coreErr.Translations().PublicMessage() {
			entryMap[lang][coreErr.UID()+"_public_msg"] = localeEntry{
				Other: translation,
			}
		}

		for lang, translation := range coreErr.Translations().InternalMessage() {
			entryMap[lang][coreErr.UID()+"_internal_msg"] = localeEntry{
				Other: translation,
			}
		}

		for argName, translations := range coreErr.Translations().Arguments() {
			for lang, translation := range translations {
				entryMap[lang][coreErr.UID()+"_argument_"+argName] = localeEntry{
					Other: translation,
				}
			}
		}

		for lang, translation := range coreErr.Translations().Deprecated() {
			entryMap[lang][coreErr.UID()+"_deprecated"] = localeEntry{
				Other: translation,
			}
		}

		for lang, translation := range coreErr.Translations().Description() {
			entryMap[lang][coreErr.UID()+"_description"] = localeEntry{
				Other: translation,
			}
		}

		for lang, translation := range coreErr.Translations().Title() {
			entryMap[lang][coreErr.UID()+"_title"] = localeEntry{
				Other: translation,
			}
		}
	}

	for lang, v := range entryMap {
		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		enc.Indent = ""
		err := enc.Encode(v)
		if err != nil {
			return fmt.Errorf("failed to encode toml: %w", err)
		}

		fileName := filepath.Join(cfg.OutputPath, fmt.Sprintf("locale.%s.toml", lang))

		err = os.WriteFile(fileName, buf.Bytes(), 0o600)
		if err != nil {
			return fmt.Errorf("failed to write %s error locale messages to file: %w", lang, err)
		}
	}

	return nil
}

func errorConstructorParams(coreErr core.Error) string {
	var res string

	for i, field := range coreErr.Arguments() {
		name := toParamName(field.Name())
		res += name + " " + typeFromArgumentType(field.Typ())

		if i != len(coreErr.Arguments())-1 {
			res += ", "
		}
	}

	return res
}

func typeFromArgumentType(argTyp core.ArgumentType) string {
	switch argTyp {
	case core.ArgumentTypeString:
		return "string"
	case core.ArgumentTypeInt:
		return "int"
	case core.ArgumentTypeFloat:
		return "float"
	case core.ArgumentTypeBool:
		return "bool"
	case core.ArgumentTypeTimestamp:
		return "time.Time"
	case core.ArgumentTypeUnknown:
		fallthrough
	default:
		panic("unknown argument type")
	}
}

func toParamName(name string) string {
	_, ok := keywords[name]
	if ok {
		return "arg" + strcase.ToCamel(name)
	}

	return name
}
