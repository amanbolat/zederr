package command

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/amanbolat/zederr/internal/codegen/core"
	"github.com/amanbolat/zederr/internal/codegen/input"
	"github.com/amanbolat/zederr/internal/codegen/output"
)

func NewGen() *cobra.Command {
	cfg := core.Config{}

	genCmd := &cobra.Command{
		Use:          "gen",
		Short:        "Generates error codes and messages.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := generateCode(cfg); err != nil {
				return err
			}

			slog.Info("successfully generated code", slog.String("go output path", cfg.ExportGo.OutputPath))

			return nil
		},
	}

	setupGenFlags(genCmd.PersistentFlags(), &cfg)

	return genCmd
}

func setupGenFlags(flagSet *pflag.FlagSet, cfg *core.Config) {
	flagSet.StringVar(&cfg.SpecPath, "spec", "./zederr_spec.yaml", "zederr specification file")
	flagSet.StringVar(&cfg.ExportGo.OutputPath, "go-out", "./gen/zederr", "output path for generated Go code")
	flagSet.StringVar(&cfg.ExportGo.PackageName, "go-pkg-name", "zederr", "package name for generated Go code")
}

func generateCode(cfg core.Config) error {
	importer := input.NewYAMLImporter()
	goExporter := output.NewGoExporter()
	manager := core.NewManager(importer, goExporter)

	err := manager.Generate(cfg)
	if err != nil {
		return err
	}

	return nil
}
