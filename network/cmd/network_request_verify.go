package cmd

import (
	"context"
	"os"

	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/numbers"
	launchtypes "github.com/ignite/network/x/launch/types"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network"
	"github.com/ignite/apps/network/network/networkchain"
	"github.com/ignite/apps/network/network/networktypes"
)

// NewNetworkRequestVerify verify the request and simulate the chain.
func NewNetworkRequestVerify() *cobra.Command {
	c := &cobra.Command{
		Use:   "verify [launch-id] [number<,...>]",
		Short: "Verify the request and simulate the chain genesis from them",
		Long: `The "verify" command applies selected requests to the genesis of a chain locally
to verify that approving these requests will result in a valid genesis that
allows a chain to launch without issues. This command does not approve requests,
only checks them.
`,
		RunE: networkRequestVerifyHandler,
		Args: cobra.ExactArgs(2),
	}

	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	return c
}

func networkRequestVerifyHandler(cmd *cobra.Command, args []string) error {
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

	// get the list of request ids
	ids, err := numbers.ParseList(args[1])
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	// verify the requests
	if err := verifyRequests(cmd.Context(), cacheStorage, nb, launchID, ids...); err != nil {
		session.Printf("%s Request(s) %s not valid\n", icons.NotOK, numbers.List(ids, "#"))
		return err
	}

	return session.Printf("%s Request(s) %s verified\n", icons.OK, numbers.List(ids, "#"))
}

// verifyRequests initializes the chain from the launch ID in a temporary directory
// and simulate the launch of the chain from genesis with the request IDs.
func verifyRequests(
	ctx context.Context,
	cacheStorage cache.Storage,
	nb NetworkBuilder,
	launchID uint64,
	requestIDs ...uint64,
) error {
	// initialize the chain for simulation
	c, n, genesisInformation, cleanup, err := initializeSimulationEnvironment(ctx, nb, launchID)
	if err != nil {
		return err
	}
	defer cleanup()

	// fetch the requests from the network
	requests, err := n.RequestFromIDs(ctx, launchID, requestIDs...)
	if err != nil {
		return err
	}

	return c.SimulateRequests(
		ctx,
		cacheStorage,
		genesisInformation,
		requests,
	)
}

// verifyRequestsFromRequestContents initializes the chain from the launch ID in a temporary directory
// and simulate the launch of the chain from genesis with the request contents.
func verifyRequestsFromRequestContents(
	ctx context.Context,
	cacheStorage cache.Storage,
	nb NetworkBuilder,
	launchID uint64,
	requestContents ...launchtypes.RequestContent,
) error {
	// initialize the chain for simulation
	c, _, genesisInformation, cleanup, err := initializeSimulationEnvironment(ctx, nb, launchID)
	if err != nil {
		return err
	}
	defer cleanup()

	return c.SimulateRequests(
		ctx,
		cacheStorage,
		genesisInformation,
		networktypes.RequestsFromRequestContents(launchID, requestContents),
	)
}

// initializeSimulationEnvironment initializes the chain from the launch ID in a temporary directory for simulating requests.
func initializeSimulationEnvironment(
	ctx context.Context,
	nb NetworkBuilder,
	launchID uint64,
) (
	c *networkchain.Chain,
	n network.Network,
	gi networktypes.GenesisInformation,
	cleanup func(),
	err error,
) {
	n, err = nb.Network()
	if err != nil {
		return c, n, gi, cleanup, err
	}

	// fetch the current genesis information and the requests for the chain for simulation
	gi, err = n.GenesisInformation(ctx, launchID)
	if err != nil {
		return c, n, gi, cleanup, err
	}

	// initialize the chain with a temporary dir
	chainLaunch, err := n.ChainLaunch(ctx, launchID)
	if err != nil {
		return c, n, gi, cleanup, err
	}

	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		return c, n, gi, cleanup, err
	}

	// TODO avoid init chain two times with the prepare command
	c, err = nb.Chain(
		networkchain.SourceLaunch(chainLaunch),
		networkchain.WithHome(homeDir),
		networkchain.WithKeyringBackend(chaincmd.KeyringBackendTest),
	)
	if err != nil {
		os.RemoveAll(homeDir)
		return c, n, gi, cleanup, err
	}

	return c, n, gi, func() { os.RemoveAll(homeDir) }, nil
}
