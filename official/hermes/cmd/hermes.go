package cmd

import (
	"github.com/spf13/cobra"
)

const (
	flagConfig = "config"
)

// NewHermes creates a new Hermes relayer command that holds
// some other sub commands related to hermes relayer.
func NewHermes() *cobra.Command {
	c := &cobra.Command{
		Use:           "hermes [command]",
		Aliases:       []string{"h"},
		Short:         "Hermes relayer wrapper",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// configure flags.
	c.PersistentFlags().StringP(flagConfig, "c", "", "set a custom Hermes config file")

	// add sub commands.
	c.AddCommand(
		NewHermesConfigure(),
		NewHermesStart(),
		NewHermesKeys(),
		NewHermesExecute(),
	)

	return c
}

func getConfig(cmd *cobra.Command) string {
	config, _ := cmd.Flags().GetString(flagConfig)
	return config
}
