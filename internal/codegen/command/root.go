package command

import (
	"github.com/spf13/cobra"
)

func NewRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "zederr",
		Short: "zederr is a tool for generating error codes and messages.",
	}

	rootCmd.AddCommand(NewGen())

	return rootCmd
}
