package cmd

import "github.com/ignite/cli/v29/ignite/services/plugin"

// GetCommands returns the list of fee-abstraction app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "fee-abstraction [command]",
			Short: "fee-abstraction is an awesome Ignite application!",
			Commands: []*plugin.Command{
				{
					Use:   "hello",
					Short: "Say hello to the world of ignite!",
				},
			},
		},
	}
}
