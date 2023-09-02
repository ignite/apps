package cmd

import (
	"github.com/spf13/cobra"
)

// NewRelayer creates a new relayer command that holds
// some other sub commands related to hermes relayer.
func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:           "hermes [command]",
		Aliases:       []string{"h"},
		Short:         "Hermes relayer wrapper",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// add sub commands.
	c.AddCommand(
		NewHermesKeys(),
		NewHermesConfigure(),
		NewHermesStart(),
		NewHermesExecute(),
	)

	return c
}
