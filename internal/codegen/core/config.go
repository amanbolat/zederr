package core

import (
	"io"
)

type Config struct {
	Source io.Reader

	GoExporterConfig GoExporterConfig
}

type GoExporterConfig struct {
	PackageName string
	Output      io.Writer
	OutputPath  string
}
