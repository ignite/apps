package cmd

import "github.com/ignite/cli/v28/ignite/services/plugin"

// GetCommands returns the list of wasm app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "wasm [command]",
			Short: "wasm is an awesome Ignite application!",
			Commands: []*plugin.Command{
				{
					Use:   "hello",
					Short: "Say hello to the world of ignite!",
				},
			},
		},
	}
}
