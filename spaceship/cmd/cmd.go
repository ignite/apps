package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of spaceship app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "spaceship [command]",
			Short: "spaceship is an awesome Ignite application!",
			Commands: []*plugin.Command{
				{
					Use:   "ssh",
					Short: "deploy your chain trough ssh",
					Commands: []*plugin.Command{
						{
							Use:   "dev",
							Short: "deploy your chain into a development mode",
						},
						{
							Use:   "deploy",
							Short: "deploy your chain into a production mode",
						},
					},
				},
			},
		},
	}
}
