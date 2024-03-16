package input

import (
	"fmt"
	"io"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"

	"github.com/amanbolat/zederr/internal/codegen/core"
)

type YAMLImporter struct{}

func NewYAMLImporter() *YAMLImporter {
	return &YAMLImporter{}
}

func (i *YAMLImporter) Import(src io.Reader) (core.Spec, error) {
	if src == nil {
		return core.Spec{}, fmt.Errorf("source is nil")
	}

	dec := yaml.NewDecoder(src)

	var yamlSpec ErrorListSpecification

	err := dec.Decode(&yamlSpec)
	if err != nil {
		return core.Spec{}, err
	}

	if len(yamlSpec.Errors) == 0 {
		return core.Spec{}, fmt.Errorf("no error entries found in the file")
	}

	defaultLocale, err := language.Parse(yamlSpec.DefaultLocale)
	if err != nil {
		return core.Spec{}, fmt.Errorf("failed to parse default locale: %w", err)
	}

	errBuilder, err := core.NewErrorBuilder(yamlSpec.SpecVersion, yamlSpec.Domain, yamlSpec.Namespace, yamlSpec.DefaultLocale)
	if err != nil {
		return core.Spec{}, fmt.Errorf("failed to create error builder: %w", err)
	}

	zedErrors := make([]core.Error, 0, len(yamlSpec.Errors))

	for _, entry := range yamlSpec.Errors {
		var args []core.Argument

		for _, rawArg := range entry.Arguments {
			arg, err := core.NewArgument(rawArg.Name, rawArg.Description, rawArg.Type)
			if err != nil {
				return core.Spec{}, err
			}

			args = append(args, arg)
		}

		localization := core.NewLocalization()

		if entry.Localization != nil {
			for _, tr := range entry.Localization.Title {
				err = localization.AddTitleTranslation(tr.Lang, tr.Value)
				if err != nil {
					return core.Spec{}, err
				}
			}

			for _, tr := range entry.Localization.Description {
				err = localization.AddDescriptionTranslation(tr.Lang, tr.Value)
				if err != nil {
					return core.Spec{}, err
				}
			}

			for _, tr := range entry.Localization.PublicMessage {
				err = localization.AddPublicMessageTranslation(tr.Lang, tr.Value)
				if err != nil {
					return core.Spec{}, err
				}
			}

			for _, tr := range entry.Localization.InternalMessage {
				err = localization.AddInternalMessageTranslation(tr.Lang, tr.Value)
				if err != nil {
					return core.Spec{}, err
				}
			}

			for _, tr := range entry.Localization.Deprecated {
				err = localization.AddDeprecatedTranslation(tr.Lang, tr.Value)
				if err != nil {
					return core.Spec{}, err
				}
			}

			for _, arg := range entry.Localization.Arguments {
				for _, tr := range arg.Description {
					err = localization.AddArgumentTranslation(arg.Name, tr.Value, tr.Lang)
					if err != nil {
						return core.Spec{}, err
					}
				}
			}
		}

		zedErr, err := errBuilder.NewError(
			entry.Code,
			entry.GRPCCode,
			entry.HTTPCode,
			entry.Description,
			entry.Title,
			entry.PublicMessage,
			entry.InternalMessage,
			entry.Deprecated,
			args,
			localization,
		)
		if err != nil {
			return core.Spec{}, err
		}

		zedErrors = append(zedErrors, zedErr)
	}

	spec := core.Spec{
		Version:       yamlSpec.SpecVersion,
		Domain:        yamlSpec.Domain,
		Namespace:     yamlSpec.Namespace,
		DefaultLocale: defaultLocale,
		Errors:        zedErrors,
	}

	return spec, nil
}
