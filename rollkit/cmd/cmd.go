package cmd

import (
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

// GetCommands returns the list of app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:     "rollkit [command]",
			Aliases: []string{"r"},
			Short:   "Ignite rollkit integration",
			Commands: []*plugin.Command{
				{
					Use:   "add",
					Short: "Add rollkit support",
					Long:  "Add rollkit support to your Cosmos SDK chain",
					Flags: []*plugin.Flag{
						{
							Name:      flagPath,
							Usage:     "path of the app",
							Shorthand: "p",
							Type:      plugin.FlagTypeString,
						},
					},
				},
				{
					Use:   "init",
					Short: "Init rollkit support",
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
