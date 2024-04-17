package main

import (
	"fmt"
	"os"

	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/spf13/cobra"

	chaininfo "github.com/ignite/apps/examples/chain-info/cmd"
	flags "github.com/ignite/apps/examples/flags/cmd"
	healthmonitor "github.com/ignite/apps/examples/health-monitor/cmd"
	helloworld "github.com/ignite/apps/examples/hello-world/cmd"
	hooks "github.com/ignite/apps/examples/hooks/cmd"
	explorer "github.com/ignite/apps/explorer/cmd"
	hermes "github.com/ignite/apps/hermes/cmd"
	marketplace "github.com/ignite/apps/marketplace/cmd"
	wasm "github.com/ignite/apps/wasm/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "apps",
		Short: "debug apps commands",
	}

	// Add apps with ignite app commands.
	newCmdFromApp(rootCmd, chaininfo.GetCommands())
	newCmdFromApp(rootCmd, flags.GetCommands())
	newCmdFromApp(rootCmd, healthmonitor.GetCommands())
	newCmdFromApp(rootCmd, helloworld.GetCommands())
	newCmdFromApp(rootCmd, hooks.GetCommands())
	newCmdFromApp(rootCmd, explorer.GetCommands())
	// Add ignite app commands for debugging here.

	// Add apps with cobra commands.
	rootCmd.AddCommand(
		hermes.NewRelayer(),
		marketplace.NewMarketplace(),
		wasm.NewWasm(),
		// Add cobra commands for debugging here.
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCmdFromApp(rootCmd *cobra.Command, commands []*plugin.Command) {
	for _, cmd := range commands {
		cobraCmd, err := cmd.ToCobraCommand()
		if err != nil {
			panic(err)
		}

		newCmdFromApp(cobraCmd, cmd.GetCommands())
		rootCmd.AddCommand(cobraCmd)
	}
}
