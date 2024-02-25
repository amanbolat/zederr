package gen

import (
	"fmt"
	"os"

	"github.com/amanbolat/zederr/internal/codegen/core"
	"github.com/amanbolat/zederr/internal/codegen/input"
	"github.com/amanbolat/zederr/internal/codegen/output"
	"github.com/urfave/cli/v2"
)

var CmdGen = cli.Command{
	Name:    "gen",
	Aliases: nil,
	Usage:   "Generate code based on errors specification.",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:     "spec",
			Required: true,
			Usage:    "USAGE",
			Aliases:  []string{"s"},
		},
		&cli.PathFlag{
			Name:     "output",
			Required: true,
			Aliases:  []string{"o"},
		},
	},
	Action: func(cmdCtx *cli.Context) error {
		return generateCode(cmdCtx)
	},
}

func generateCode(cmdCtx *cli.Context) error {
	specFilePath := cmdCtx.Path("spec")

	importer := input.NewYAMLImporter()
	goExporter := output.NewGoExporter()
	manager := core.NewManager(importer, goExporter)

	specFile, err := os.Open(specFilePath)
	if err != nil {
		return fmt.Errorf("failed to open spec file: %w", err)
	}
	defer func() {
		_ = specFile.Close()
	}()

	outputPath := cmdCtx.Path("output")
	if outputPath == "" {
		return fmt.Errorf("output path was not provided")
	}

	err = os.MkdirAll(outputPath, 0750)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	err = manager.Generate(core.Config{
		Source: specFile,
		GoExporterConfig: core.GoExporterConfig{
			PackageName: "zederr",
			Output:      nil,
			OutputPath:  outputPath,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
