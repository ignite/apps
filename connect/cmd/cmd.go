package cmd

import (
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// GetCommands returns the list of app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:     "connect [command]",
			Aliases: []string{"c"},
			Short:   "Interact with any Cosmos SDK based blockchain using Ignite Connect",
			Long:    "Ignite Connect is an app that allows you to interact with any Cosmos SDK based blockchain.\n It leverages AutoCLI from client/v2.",
			Commands: []*plugin.Command{
				{
					Use:   "discover",
					Short: "Discover chains to connect to",
				},
				{
					Use:     "add [chain]",
					Short:   "Add a chain to interact with",
					Aliases: []string{"a", "init"},
				},
				{
					Use:   "version",
					Short: "Print the version of the Ignite Connect app",
				},
			},
		},
	}
}
