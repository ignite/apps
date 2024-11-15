package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of web app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "web [command]",
			Short: "Ignite chain dashboard",
			Commands: []*plugin.Command{
				{
					Use:   "add",
					Short: "Add the chain-admin dashboard on your app",
					Long:  "Scaffold the chain-admin dashboard, and easily deploy a frontend using Next.js and Cosmos Kit",
				},
			},
		},
	}
}
