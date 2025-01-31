package cmd

import (
	"fmt"

	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// GetCommands returns the list of app commands.
func GetCommands(availableChains []string) []*plugin.Command {
	cmd := []*plugin.Command{
		{
			Use:     "connect [command]",
			Aliases: []string{"c"},
			Short:   "Interact with any Cosmos SDK based blockchain using Ignite Connect",
			Long:    "Connect allows you to interact with any Cosmos SDK based blockchain.",
			Commands: []*plugin.Command{
				{
					Use:   "discover",
					Short: "Discover chains to connect to",
				},
				{
					Use:     "add <chain> [endpoint]",
					Short:   "Add a chain to interact with",
					Long:    "Add a chain to interact with. If a chain and endpoint are provided, the chain will be added without prompting",
					Aliases: []string{"init"},
				},
				{
					Use:     "remove <chain>",
					Short:   "Remove a chain from Connect",
					Aliases: []string{"rm"},
				},
				{
					Use:   "version",
					Short: "Display Connect version",
				},
			},
		},
	}

	for _, name := range availableChains {
		cmd[0].Commands = append(cmd[0].Commands,
			&plugin.Command{
				Use:   name,
				Short: fmt.Sprintf("Commands for %s chain", name),
			},
		)
	}

	return cmd
}
