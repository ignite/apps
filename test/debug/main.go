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
)

var rootCmd = &cobra.Command{
	Use:   "apps",
	Short: "debug apps commands",
}

func newCmdFromApp(name string, cmds []*plugin.Command) *cobra.Command {
	newCmd := &cobra.Command{Use: fmt.Sprintf("%s [app]", name)}
	for _, cmd := range cmds {
		cobraCmd, err := cmd.ToCobraCommand()
		if err != nil {
			panic(err)
		}
		newCmd.AddCommand(cobraCmd)
	}
	return newCmd
}

func main() {
	rootCmd.AddCommand(
		marketplace.NewMarketplace(),
		hermes.NewHermes(),
		newCmdFromApp("explorer", explorer.GetCommands()),
		newCmdFromApp("chain-info", chaininfo.GetCommands()),
		newCmdFromApp("flags", flags.GetCommands()),
		newCmdFromApp("health-monitor", healthmonitor.GetCommands()),
		newCmdFromApp("hello-world", helloworld.GetCommands()),
		newCmdFromApp("hooks", hooks.GetCommands()),
		// Add commands for debugging here.
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
