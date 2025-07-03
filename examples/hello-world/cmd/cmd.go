package cmd

import "github.com/ignite/cli/v29/ignite/services/plugin"

// GetCommands returns the list of hello-world app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "hello-world",
			Short: "Say hello to the world of ignite!",
		},
	}
}
