package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of web app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			PlaceCommandUnder: "scaffold",
			Use:               "cca",
			Short:             "Ignite CCA scaffolds a Cosmos SDK chain frontend using a `create-cosmos-app` template",
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
