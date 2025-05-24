package cmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkReward creates a new chain reward command.
func NewNetworkReward() *cobra.Command {
	c := &cobra.Command{
		Use:    "reward",
		Short:  "Manage network rewards",
		Hidden: true,
	}
	c.AddCommand(
		NewNetworkRewardSet(),
		NewNetworkRewardRelease(),
	)
	return c
}
