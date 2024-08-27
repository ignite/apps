package cmd

import (
	"context"
	"fmt"

	"github.com/ignite/cli/v29/ignite/services/plugin"
)

// ExecuteHello executes the hello subcommand.
func ExecuteHello(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	fmt.Println("Hello, world!")
	return nil
}
