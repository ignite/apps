package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of web app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "web [command]",
			Short: "web is an awesome Ignite application!",
			Commands: []*plugin.Command{
				{
					Use:   "hello",
					Short: "Say hello to the world of ignite!",
				},
			},
		},
	}
}
