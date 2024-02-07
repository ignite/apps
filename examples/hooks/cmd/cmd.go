package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of hooks app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "hooks",
			Short: "This is an example Ignite App that demonstrates usage of hooks",
			Long:  `To use either run "ignite chain build" or "ignite chain serve" and see the output.`,
		},
	}
}
