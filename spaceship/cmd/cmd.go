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
					Use:   "deploy [host]",
					Short: "deploy the chain",
					Flags: append(defaultFlags,
						&plugin.Flag{
							Name:      flagInitChain,
							Shorthand: "i",
							Usage:     "run init chain and create the home folder",
							Type:      plugin.FlagTypeBool,
						},
						&plugin.Flag{
							Name:         flagFaucet,
							Shorthand:    "f",
							Usage:        "create a chain faucet",
							Type:         plugin.FlagTypeBool,
							DefaultValue: "false",
						},
						&plugin.Flag{
							Name:         flagFaucetPort,
							Usage:        "chain faucet port",
							Type:         plugin.FlagTypeUint64,
							DefaultValue: "8009",
						},
					),
				},
				{
					Use:   "log [host]",
					Short: "get remote logs",
					Commands: []*plugin.Command{
						{
							Use:   "chain",
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
							Use:   "faucet",
							Short: "get faucet logs if its running",
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
					},
				},
				{
					Use:   "status [host]",
					Short: "get chain status if its running",
					Flags: defaultFlags,
				},
				{
					Use:   "restart [host]",
					Short: "restart the chain",
					Flags: defaultFlags,
				},
				{
					Use:   "stop [host]",
					Short: "stop the chain",
					Flags: defaultFlags,
				},
				{
					Use:   "faucet",
					Short: "faucet commands",
					Commands: []*plugin.Command{
						{
							Use:   "status [host]",
							Short: "get faucet status if its running",
							Flags: defaultFlags,
						},
						{
							Use:   "start [host]",
							Short: "start the faucet",
							Flags: append(defaultFlags,
								&plugin.Flag{
									Name:         flagFaucetPort,
									Usage:        "chain faucet port",
									Type:         plugin.FlagTypeUint64,
									DefaultValue: "8009",
								},
							),
						},
						{
							Use:   "restart [host]",
							Short: "restart the faucet",
							Flags: append(defaultFlags,
								&plugin.Flag{
									Name:         flagFaucetPort,
									Usage:        "chain faucet port",
									Type:         plugin.FlagTypeUint64,
									DefaultValue: "8009",
								},
							),
						},
						{
							Use:   "stop",
							Short: "stop the faucet",
							Flags: defaultFlags,
						},
					},
				},
			},
		},
	}
}
