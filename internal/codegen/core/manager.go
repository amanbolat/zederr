package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type Importer interface {
	Import(reader io.Reader) (Spec, error)
}

// GoExporter is responsible for exporting the parsed errors to Go code.
type GoExporter interface {
	Export(cfg ExportGo, spec Spec) error
}

type Manager struct {
	importer   Importer
	goExporter GoExporter
}

func NewManager(importer Importer, goExporter GoExporter) *Manager {
	return &Manager{
		importer:   importer,
		goExporter: goExporter,
	}
}

func (m *Manager) Generate(cfg Config) error {
	b, err := os.ReadFile(cfg.SpecPath)
	if err != nil {
		return fmt.Errorf("failed to read spec file: %w", err)
	}

	spec, err := m.importer.Import(bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	err = m.goExporter.Export(cfg.ExportGo, spec)
	if err != nil {
		return err
	}

	return nil
}
