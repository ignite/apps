package cmd

import (
	"context"
	"fmt"

	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// ExecuteHello executes the hello subcommand.
func ExecuteHello(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}

	name, err := getNameFlag(flags)
	if err != nil {
		return err
	}
	fmt.Printf("Hello, %s!\n", name)
	return nil
}
