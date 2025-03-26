package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of spinner app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "spinner",
			Short: "App spinner example",
		},
	}
}
