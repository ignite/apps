package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of health-monitor app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "health-monitor [command]",
			Short: "health-monitor is an awesome Ignite application!",
			Commands: []*plugin.Command{
				{
					Use:   "monitor",
					Short: "monitor and print out health status of a running chain",
					Flags: []*plugin.Flag{
						{
							Name:         "json",
							Usage:        "output as JSON",
							DefaultValue: "false",
							Type:         plugin.FlagTypeBool,
						},
						{
							Name:         "path",
							Usage:        "path of the app",
							DefaultValue: ".",
							Type:         plugin.FlagTypeString,
						},
						{
							Name:         "refresh-duration",
							Shorthand:    "r",
							Usage:        "refresh duration of the monitor",
							DefaultValue: "5s",
							Type:         plugin.FlagTypeString,
						},
						{
							Name:  "rpc-address",
							Usage: "RPC address of the chain to monitor (default: current chain RPC address)",
							Type:  plugin.FlagTypeString,
						},
					},
				},
			},
		},
	}
}
