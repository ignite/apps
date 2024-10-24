package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

var defaultFlags = []*plugin.Flag{
	{
		Name:      flagUser,
		Shorthand: "u",
		Usage:     "ssh user",
		Type:      plugin.FlagTypeString,
	},
	{
		Name:      flagPort,
		Shorthand: "p",
		Usage:     "ssh port",
		Type:      plugin.FlagTypeString,
	},
	{
		Name:  flagPassword,
		Usage: "ssh user password",
		Type:  plugin.FlagTypeString,
	},
	{
		Name:      flagKey,
		Shorthand: "k",
		Usage:     "ssh key",
		Type:      plugin.FlagTypeString,
	},
	{
		Name:      flagRawKey,
		Shorthand: "r",
		Usage:     "ssh raw key",
		Type:      plugin.FlagTypeString,
	},
	{
		Name:  flagKeyPassword,
		Usage: "ssh key password",
		Type:  plugin.FlagTypeString,
	},
}

// GetCommands returns the list of spaceship app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "spaceship [command]",
			Short: "Deploy a chain remote through SSH using ignite build system",
			Commands: []*plugin.Command{
				{
					Use:   "deploy",
					Short: "deploy your chain",
					Flags: append(defaultFlags,
						&plugin.Flag{
							Name:      flagInitChain,
							Shorthand: "i",
							Usage:     "run init chain and create the home folder",
							Type:      plugin.FlagTypeBool,
						},
					),
				},
				{
					Use:   "log",
					Short: "get chain logs if its running",
					Flags: append(defaultFlags,
						&plugin.Flag{
							Name:         flagLines,
							Shorthand:    "l",
							Usage:        "number of lines of chain logs",
							Type:         plugin.FlagTypeInt,
							DefaultValue: "100",
						},
						&plugin.Flag{
							Name:  flagRealTime,
							Usage: "show the logs in the real time",
							Type:  plugin.FlagTypeBool,
						},
					),
				},
				{
					Use:   "status",
					Short: "get chain status if its running",
					Flags: defaultFlags,
				},
				{
					Use:   "restart",
					Short: "restart your chain",
					Flags: defaultFlags,
				},
				{
					Use:   "stop",
					Short: "stop your chain",
					Flags: defaultFlags,
				},
			},
		},
	}
}
