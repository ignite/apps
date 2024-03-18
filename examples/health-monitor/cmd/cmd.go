package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

const (
	flagJSON            = "json"
	flagPath            = "path"
	flagRefreshDuration = "refresh-duration"
	flagRPCAddress      = "rpc-address"
)

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
							Name:         flagJSON,
							Usage:        "output as JSON",
							DefaultValue: "false",
							Type:         plugin.FlagTypeBool,
						},
						{
							Name:         flagPath,
							Usage:        "path of the app",
							DefaultValue: ".",
							Type:         plugin.FlagTypeString,
						},
						{
							Name:         flagRefreshDuration,
							Shorthand:    "r",
							Usage:        "refresh duration of the monitor",
							DefaultValue: "5s",
							Type:         plugin.FlagTypeString,
						},
						{
							Name:  flagRPCAddress,
							Usage: "RPC address of the chain to monitor (default: current chain RPC address)",
							Type:  plugin.FlagTypeString,
						},
					},
				},
			},
		},
	}
}
