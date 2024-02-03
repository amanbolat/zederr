package main

import (
	"log/slog"
	"os"

	"github.com/amanbolat/zederr/internal/core"
	"github.com/amanbolat/zederr/internal/input"
	"github.com/amanbolat/zederr/internal/output"
	"github.com/amanbolat/zederr/internal/parser"
)

func main() {
	p := parser.NewParser("{{", "}}", nil)
	errBuilder := core.NewErrorBuilder(p)
	yamlImporter := input.NewYAMLImporter(errBuilder)
	goExporter := output.NewGoExporter()
	manager := core.NewManager(yamlImporter, goExporter)

	err := manager.Generate(core.Config{})
	if err != nil {
		slog.Error("failed to generate", slog.Any("error", err))
		os.Exit(1)
	}
}
