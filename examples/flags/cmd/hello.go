package cmd

import (
	"context"
	"fmt"

	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// ExecuteHello executes the hello subcommand.
func ExecuteHello(_ context.Context, cmd *plugin.ExecutedCommand) error {
	flags := plugin.Flags(cmd.Flags)

	name, err := flags.GetString(flagName)
	if err != nil {
		return err
	}
	fmt.Printf("Hello, %s!\n", name)
	return nil
}
