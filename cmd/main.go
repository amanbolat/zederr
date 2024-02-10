package main

import (
	"log/slog"
	"os"

	"github.com/amanbolat/zederr/pkg/codegen/core"
	"github.com/amanbolat/zederr/pkg/codegen/input"
	"github.com/amanbolat/zederr/pkg/codegen/output"
	"github.com/amanbolat/zederr/pkg/codegen/parser"
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
