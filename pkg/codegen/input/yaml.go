package input

import (
	"fmt"
	"io"

	"github.com/amanbolat/zederr/pkg/codegen/core"
	"gopkg.in/yaml.v3"
)

type YAMLImporter struct {
	errBuilder *core.ErrorBuilder
}

func NewYAMLImporter(errBuilder *core.ErrorBuilder) *YAMLImporter {
	return &YAMLImporter{
		errBuilder: errBuilder,
	}
}

func (i *YAMLImporter) Import(src io.Reader) ([]core.Error, error) {
	if src == nil {
		return nil, fmt.Errorf("source is nil")
	}

	dec := yaml.NewDecoder(src)
	var entries ErrorEntries

	err := dec.Decode(&entries)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no error entries found in the file")
	}

	var zedErrors []core.Error
	for _, entry := range entries {
		zedErr, err := i.errBuilder.NewError(entry.ID, entry.GRPCCode, entry.HTTPCode, entry.Description, entry.Translations)
		if err != nil {
			return nil, err
		}

		zedErrors = append(zedErrors, zedErr)
	}

	return zedErrors, nil
}
