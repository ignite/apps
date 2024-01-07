package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewHooks creates a new hooks command.
func NewHooks() *cobra.Command {
	c := &cobra.Command{
		Use:           "hooks",
		Short:         "This is a example Ignite App that demonstrates hooks",
		Long:          "To use either run \"ignite scaffold chain\" or \"ignite chain serve\" and see the output.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.Long)
			return nil
		},
	}

	return c
}
