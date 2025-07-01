package cmd

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/pkg/ssh"
)

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
		Name:  flagUserPassword,
		Usage: "ssh user password",
		Type:  plugin.FlagTypeString,
	},
	{
		Name:  flagPassword,
		Usage: "ask the ssh user password",
		Type:  plugin.FlagTypeBool,
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
							Usage:        "run the chain faucet",
							Type:         plugin.FlagTypeBool,
							DefaultValue: "false",
						},
						&plugin.Flag{
							Name:  flagFaucetPort,
							Usage: "chain faucet port",
							Type:  plugin.FlagTypeUint64,
						},
					),
				},
				{
					Use:   "log [host]",
					Short: "get remote logs",
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
						&plugin.Flag{
							Name:      flagAppLog,
							Shorthand: "a",
							Usage: fmt.Sprintf(
								"the app to show the log (%s)",
								strings.Join(ssh.LogTypes(), ","),
							),
							Type:         plugin.FlagTypeString,
							DefaultValue: ssh.LogChain.String(),
						},
					),
				},
				{
					Use:   "status [host]",
					Short: "get chain status if its running",
					Flags: append(defaultFlags,
						&plugin.Flag{
							Name:         flagFaucet,
							Shorthand:    "f",
							Usage:        "show faucet status",
							Type:         plugin.FlagTypeBool,
							DefaultValue: "false",
						},
					),
				},
				{
					Use:   "restart [host]",
					Short: "restart the chain",
					Flags: append(defaultFlags,
						&plugin.Flag{
							Name:         flagFaucet,
							Shorthand:    "f",
							Usage:        "run the chain faucet",
							Type:         plugin.FlagTypeBool,
							DefaultValue: "false",
						},
						&plugin.Flag{
							Name:  flagFaucetPort,
							Usage: "chain faucet port",
							Type:  plugin.FlagTypeUint64,
						},
					),
				},
				{
					Use:   "stop [host]",
					Short: "stop the chain",
					Flags: append(defaultFlags,
						&plugin.Flag{
							Name:         flagFaucet,
							Shorthand:    "f",
							Usage:        "stop the chain faucet",
							Type:         plugin.FlagTypeBool,
							DefaultValue: "false",
						},
					),
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
									Name:  flagFaucetPort,
									Usage: "chain faucet port",
									Type:  plugin.FlagTypeUint64,
								},
							),
						},
						{
							Use:   "restart [host]",
							Short: "restart the faucet",
							Flags: append(defaultFlags,
								&plugin.Flag{
									Name:  flagFaucetPort,
									Usage: "chain faucet port",
									Type:  plugin.FlagTypeUint64,
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
