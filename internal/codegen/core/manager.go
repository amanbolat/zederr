package core

import (
	"io"
)

type Importer interface {
	Import(reader io.Reader) (Spec, error)
}

// GoExporter is responsible for exporting the parsed errors to Go code.
type GoExporter interface {
	Export(cfg GoExporterConfig, spec Spec) error
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
	spec, err := m.importer.Import(cfg.Source)
	if err != nil {
		return err
	}

	err = m.goExporter.Export(cfg.GoExporterConfig, spec)
	if err != nil {
		return err
	}

	return nil
}
