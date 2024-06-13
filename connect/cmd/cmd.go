package cmd

import "github.com/spf13/cobra"

// NewConnect creates a new connect command that holds
// some other sub commands related to Ignite Connect.
func NewConnect() *cobra.Command {
	c := &cobra.Command{
		Use:           "connect [command]",
		Aliases:       []string{"c"},
		Short:         "Interact with any Cosmos SDK based blockchain using Ignite Connect",
		Long:          "Ignite Connect is an app that allows you to interact with any Cosmos SDK based blockchain.\n It leverages AutoCLI from client/v2 and is inspired by the Hubl tool",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// add sub commands.
	c.AddCommand()
	return c
}
