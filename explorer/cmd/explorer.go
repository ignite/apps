package cmd

import (
	"github.com/spf13/cobra"
)

// NewExplorer creates a new explorer command that holds
// some other sub commands related to running chain explorers like gex and etc.
func NewExplorer() *cobra.Command {
	c := &cobra.Command{
		Use:           "explorer [command]",
		Aliases:       []string{"e"},
		Short:         "Run chain explorer commands",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// add sub commands.
	c.AddCommand(
		NewGex(),
	)

	return c
}
