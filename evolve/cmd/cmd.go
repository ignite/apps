package cmd

import (
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

// GetCommands returns the list of app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:     "evolve [command]",
			Aliases: []string{"r", "e", "rollkit"},
			Short:   "Ignite evolve integration",
			Commands: []*plugin.Command{
				{
					Use:   "add",
					Short: "Add evolve support",
					Long:  "Add evolve support to your Cosmos SDK chain",
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
				{
					Use:   "add-migrate",
					Short: "Add evolve migration support",
					Long:  "Add evolve migration helpers and modules for migrating from CometBFT",
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
				{
					Use:   "init",
					Short: "Init evolve support",
					Long:  "Initialize the chain and a ev-node sequencer via ev-abci.",
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
				{
					Use:   "edit-genesis",
					Short: "Edit the genesis file to add the evolve sequencer",
					Long:  "Edit the genesis file to add the evolve sequencer module. This is an alternative to the `ignite evolve init` command, where a chain is already initialized (but not yet started).",
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
