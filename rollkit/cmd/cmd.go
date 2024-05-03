package cmd

import "github.com/spf13/cobra"

// NewRollkit creates a new rollkit command that holds
// some other sub commands related to Rollkit.
func NewRollkit() *cobra.Command {
	c := &cobra.Command{
		Use:           "rollkit [command]",
		Aliases:       []string{"r"},
		Short:         "Ignite rollkit integration",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// add sub commands.
	c.AddCommand(
		NewRollkitAdd(),
	)
	return c
}
