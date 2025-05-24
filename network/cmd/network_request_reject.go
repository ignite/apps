package cmd

import (
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/numbers"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network"
)

// NewNetworkRequestReject creates a new request reject
// command to reject requests for a chain.
func NewNetworkRequestReject() *cobra.Command {
	c := &cobra.Command{
		Use:     "reject [launch-id] [number<,...>]",
		Aliases: []string{"accept"},
		Short:   "Reject requests",
		Long: `The "reject" command is used by a chain's coordinator to reject requests.

	ignite network request reject 42 1,2,3-6,7,8

The syntax of the "reject" command is similar to that of the "approve" command.
`,
		RunE: networkRequestRejectHandler,
		Args: cobra.ExactArgs(2),
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	return c
}

func networkRequestRejectHandler(cmd *cobra.Command, args []string) error {
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

	// Get the list of request ids
	ids, err := numbers.ParseList(args[1])
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	// Submit the rejected requests
	reviewals := make([]network.Reviewal, 0)
	for _, id := range ids {
		reviewals = append(reviewals, network.RejectRequest(id))
	}
	if err := n.SubmitRequestReviewals(cmd.Context(), launchID, reviewals...); err != nil {
		return err
	}

	return session.Printf("%s Request(s) %s rejected\n", icons.OK, numbers.List(ids, "#"))
}
