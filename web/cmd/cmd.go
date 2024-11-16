package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of web app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			PlaceCommandUnder: "scaffold",
			Use:               "web",
			Short:             "Ignite chain-admin dashboard",
			Long:              "Scaffolds the chain-admin dashboard, an easy to use frontend powered by Next.js and Cosmos Kit",
			Flags: []*plugin.Flag{
				{
					Name:      flagPath,
					Usage:     "path of the app",
					Shorthand: "p",
					Type:      plugin.FlagTypeString,
				},
			},
		},
	}
}
