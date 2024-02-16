package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of explorer app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:     "explorer [command]",
			Short:   "Run chain explorer commands",
			Aliases: []string{"e"},
			Commands: []*plugin.Command{
				{
					Use:     "gex [rpc_url]",
					Short:   "Run gex",
					Aliases: []string{"g"},
				},
			},
		},
	}
}
