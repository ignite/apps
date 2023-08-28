package cmd

import (
	"github.com/spf13/cobra"
)

const (
	flagHome = "home"
	flagFrom = "from"
)

// NewHermes creates a new hermes command that holds
// some other sub commands related to hermes relayer.
func NewHermes() *cobra.Command {
	c := &cobra.Command{
		Use:           "hermes [command]",
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
		NewHermesConfigure(),
	)

	return c
}
