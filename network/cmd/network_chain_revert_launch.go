package cmd

import (
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network"
)

// NewNetworkChainRevertLaunch creates a new chain revert launch command
// to revert a launched chain.
func NewNetworkChainRevertLaunch() *cobra.Command {
	c := &cobra.Command{
		Use:   "revert-launch [launch-id]",
		Short: "Revert launch of a network as a coordinator",
		Long: `The revert launch command reverts the previously scheduled launch of a chain.

Only the coordinator of the chain can execute the launch command.

	ignite network chain revert-launch 42

After the revert launch command is executed, changes to the genesis of the chain
are allowed again. For example, validators will be able to request to join the
chain. Revert launch also resets the launch time.
`,
		Args: cobra.ExactArgs(1),
		RunE: networkChainRevertLaunchHandler,
	}

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())

	return c
}

func networkChainRevertLaunchHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	return n.RevertLaunch(cmd.Context(), launchID)
}
