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
					Long:  "Initialize the chain and add rollkit sequencer",
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
					Use:   "edit-genesis",
					Short: "Edit the genesis file to add the rollkit sequencer",
					Long:  "Edit the genesis file to add the rollkit sequencer module. This is an alternative to the `ignite rollkit init` command, where a chain is already initialized (but not yet started).",
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
					Use:     "migrate-from-cometbft",
					Aliases: []string{"migrate-from-comet"},
					Short:   "Migrate from CometBFT to Rollkit",
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
