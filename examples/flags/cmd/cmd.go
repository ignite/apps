package cmd

import (
	"fmt"

	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/spf13/pflag"
)

// GetCommands returns the list of flags app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "flags [command]",
			Short: "flsgs is a simple application that demonstrates use of cli flags and args in Ignite applications",
			Commands: []*plugin.Command{
				{
					Use:   "hello",
					Short: "Say hello to the user!",
				},
				{
					Use:   "cowsay",
					Short: "Cow says hello to the user!",
					Flags: []*plugin.Flag{
						{
							Name:         "type",
							Shorthand:    "t",
							Usage:        "Type of the cow! (Try cheese)",
							DefaultValue: "default",
							Type:         plugin.FlagTypeString,
						},
					},
				},
			},
			Flags: []*plugin.Flag{
				{
					Name:         "name",
					Shorthand:    "n",
					Usage:        "Name of the one you want to say hello to!",
					DefaultValue: "Ignite",
					Type:         plugin.FlagTypeString,
					Persistent:   true,
				},
			},
		},
	}
}

func getNameFlag(flags *pflag.FlagSet) (string, error) {
	name, err := flags.GetString("name")
	if err != nil {
		return "", fmt.Errorf("could not get --name flag: %w", err)
	}

	if name == "" {
		return "Ignite", nil
	}
	return name, nil
}
