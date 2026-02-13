package cmd

import (
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

// GetCommands returns the list of app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:     "gnovm [command]",
			Aliases: []string{"gno"},
			Short:   "Ignite GnoVM integration",
			Commands: []*plugin.Command{
				{
					Use:   "add",
					Short: "Add GnoVM support",
					Long:  "Add GnoVM support to your Cosmos SDK chain",
					Flags: []*plugin.Flag{
						{
							Name:         flagPath,
							Usage:        "path of the app",
							Shorthand:    "p",
							Type:         plugin.FlagTypeString,
							DefaultValue: ".",
						},
					},
				},
			},
		},
	}
}
