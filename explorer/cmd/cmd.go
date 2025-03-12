package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

const flagRPCAddress = "rpc-address"

// GetCommands returns the list of explorer app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:     "explorer [command]",
			Short:   "Run chain explorer commands",
			Aliases: []string{"e"},
			Commands: []*plugin.Command{
				{
					Use:     "gex",
					Short:   "Run gex explorer",
					Aliases: []string{"g"},
					Flags: []*plugin.Flag{
						{
							Name:         flagRPCAddress,
							Usage:        "The chain RPC address",
							DefaultValue: "http://localhost:26657",
							Type:         plugin.FlagTypeString,
						},
					},
				},
				{
					Use:     "pingpub",
					Short:   "Run Ping pub explorer",
					Aliases: []string{"p"},
					Flags: []*plugin.Flag{
						{
							Name:      flagPath,
							Usage:     "path of the app",
							Shorthand: "p",
							Type:      plugin.FlagTypeString,
						},
						{
							Name:         flagPort,
							Usage:        "port to run the explorer on",
							Type:         plugin.FlagTypeUint,
							DefaultValue: "8080",
						},
					},
				},
			},
		},
	}
}
