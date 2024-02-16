package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// ExecuteInfo executes the info subcommand.
func ExecuteInfo(ctx context.Context, cmd *plugin.ExecutedCommand, c *chain.Chain) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	write := func(s string, v interface{}) {
		fmt.Fprintf(w, "%s:\t%v\n", s, v)
	}

	write("Version", c.Version)
	write("App Path", c.AppPath())
	write("Config Path", c.ConfigPath())
	init, err := c.IsInitialized()
	if err != nil {
		return fmt.Errorf("could not find out if the chain is initialized: %w", err)
	}
	write("Is Initialized", init)
	bin, err := c.Binary()
	if err != nil {
		return fmt.Errorf("could not find out chain's binary file name: %w", err)
	}
	write("Binary File", bin)
	w.Flush()

	return nil
}
