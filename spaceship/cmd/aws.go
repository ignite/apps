package cmd

import (
	"context"
	"fmt"

	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// ExecuteAWS executes the aws deploy subcommand.
func ExecuteAWS(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	fmt.Println("Hello, world!")
	return nil
}
