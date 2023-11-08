package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	explorer "github.com/ignite/apps/explorer/cmd"
	hermes "github.com/ignite/apps/hermes/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "apps",
	Short: "debug apps commands",
}

func main() {
	rootCmd.AddCommand(
		explorer.NewExplorer(),
		hermes.NewRelayer(),
		// Add commands for debugging here.
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
