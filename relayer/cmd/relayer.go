package cmd

import (
	"github.com/spf13/cobra"
)

const (
	flagHome = "home"
	flagFrom = "from"
)

// NewRelayer creates a new relayer command that holds
// some other sub commands related to hermes relayer.
func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:           "relayer [command]",
		Aliases:       []string{"h"},
		Short:         "",
		Long:          ``,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// configure flags.
	//c.PersistentFlags().BoolVar(&local, flagLocal, false, "blabla")

	// add sub commands.
	c.AddCommand(
		NewRelayerConfigure(),
	)

	return c
}
