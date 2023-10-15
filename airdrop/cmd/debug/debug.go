package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ignite/apps/airdrop/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "airdrop",
	Short: "debug command for CLI airdrop plugin",
}

func main() {
	rootCmd.AddCommand(cmd.NewAirdrop())
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
