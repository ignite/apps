package cmd

import "github.com/ignite/cli/v29/ignite/services/plugin"

const flagOutput = "output"

// GetCommands returns the list of chain-info app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "chain-info",
			Short: "chain-info is a simple application that demonstrates how to get information or manipulate the chains",
			Commands: []*plugin.Command{
				{
					Use:   "info",
					Short: "Prints out some basic informations about the chain in the current directory",
				},
				{
					Use:   "build",
					Short: "Builds the chain app in the current directory with help of ignite helper functions",
					Flags: []*plugin.Flag{
						{
							Name:         flagOutput,
							Shorthand:    "o",
							Usage:        "The path to output binary file",
							DefaultValue: ".",
							Type:         plugin.FlagTypeString,
						},
					},
				},
			},
		},
	}
}
