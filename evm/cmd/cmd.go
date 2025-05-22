package cmd

import (
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// GetCommands returns the list of app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:     "evm [command]",
			Aliases: []string{"e"},
			Short:   "Ignite EVM integration",
			Commands: []*plugin.Command{
				{
					Use:   "add",
					Short: "Add EVM support",
					Long:  "Add EVM support to your Cosmos SDK chain",
					Flags: []*plugin.Flag{
						{
							Name:      flagPath,
							Usage:     "path of the app",
							Shorthand: "p",
							Type:      plugin.FlagTypeString,
						},
					},
				},
			},
		},
	}
}
