package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/gitpod"
	"github.com/ignite/cli/v28/ignite/pkg/goenv"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network"
	"github.com/ignite/apps/network/network/networkchain"
	"github.com/ignite/apps/network/network/networktypes"
)

const (
	flagForce = "force"
)

// NewNetworkChainPrepare returns a new command to prepare the chain for launch.
func NewNetworkChainPrepare() *cobra.Command {
	c := &cobra.Command{
		Use:   "prepare [launch-id]",
		Short: "Prepare the chain for launch",
		Long: `The prepare command prepares a validator node for the chain launch by generating
the final genesis and adding IP addresses of peers to the validator's
configuration file.

	ignite network chain prepare 42

By default, Ignite uses "$HOME/spn/LAUNCH_ID" as the data directory. If you used
a different data directory when initializing the node, use the "--home" flag and
set the correct path to the data directory.

Ignite generates the genesis file in "config/genesis.json" and adds peer IPs by
modifying "config/config.toml".

The prepare command should be executed after the coordinator has triggered the
chain launch and finalized the genesis with "ignite network chain launch". You
can force Ignite to run the prepare command without checking if the launch has
been triggered with the "--force" flag (this is not recommended).

After the prepare command is executed the node is ready to be started.
`,
		Args: cobra.ExactArgs(1),
		RunE: networkChainPrepareHandler,
	}

	flagSetClearCache(c)
	c.Flags().BoolP(flagForce, "f", false, "force the prepare command to run even if the chain is not launched")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetCheckDependencies())

	return c
}

func networkChainPrepareHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	force, _ := cmd.Flags().GetBool(flagForce)

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

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

	// fetch chain information
	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	if !force && chainLaunch.Metadata.Cli.Version != "" && !chainLaunch.Metadata.IsCurrentVersion() {
		return fmt.Errorf(`chain %d has been published with a different version of the plugin (%s, current version is %s)
this may result in a genesis that is different from other validators' genesis
use --force to prepare anyway`,
			launchID,
			chainLaunch.Metadata.Cli.Version,
			networktypes.Version,
		)
	}

	if !force && !chainLaunch.LaunchTriggered {
		return fmt.Errorf("chain %d launch has not been triggered yet. use --force to prepare anyway", launchID)
	}

	networkOptions := []networkchain.Option{
		networkchain.WithKeyringBackend(chaincmd.KeyringBackendTest),
	}

	if flagGetCheckDependencies(cmd) {
		networkOptions = append(networkOptions, networkchain.CheckDependencies())
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch), networkOptions...)
	if err != nil {
		return err
	}

	if err := prepareFromGenesisInformation(
		cmd,
		cacheStorage,
		launchID,
		n,
		c,
		chainLaunch,
	); err != nil {
		return err
	}

	chainHome, err := c.Home()
	if err != nil {
		return err
	}
	binaryName, err := c.BinaryName()
	if err != nil {
		return err
	}
	binaryDir := filepath.Dir(filepath.Join(goenv.Bin(), binaryName))

	session.Printf("%s Chain is prepared for launch\n", icons.OK)
	session.Println("\nYou can start your node by running the following command:")
	commandStr := fmt.Sprintf("%s start --home %s", binaryName, chainHome)
	if gitpod.IsOnGitpod() {
		// Gitpod requires to enable proxy-tunnel tool
		commandStr = fmt.Sprintf(
			"ignite network tool proxy-tunnel %s/spn.yml & %s",
			chainHome, commandStr,
		)
	}
	session.Printf("\t%s/%s\n", binaryDir, colors.Info(commandStr))

	return nil
}

// prepareFromGenesisInformation prepares the genesis of the chain from the queried genesis information from the launch ID of the chain.
func prepareFromGenesisInformation(
	cmd *cobra.Command,
	cacheStorage cache.Storage,
	launchID uint64,
	n network.Network,
	c *networkchain.Chain,
	chainLaunch networktypes.ChainLaunch,
) error {
	var (
		rewardsInfo           networktypes.Reward
		lastBlockHeight       int64
		consumerUnbondingTime int64
	)

	// fetch the information to construct genesis
	genesisInformation, err := n.GenesisInformation(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	// fetch the info for rewards if the consumer revision height is defined
	if chainLaunch.ConsumerRevisionHeight > 0 {
		rewardsInfo, lastBlockHeight, consumerUnbondingTime, err = n.RewardsInfo(
			cmd.Context(),
			launchID,
			chainLaunch.ConsumerRevisionHeight,
		)
		if err != nil {
			return err
		}
	}

	spnChainID, err := n.ChainID(cmd.Context())
	if err != nil {
		return err
	}

	return c.Prepare(
		cmd.Context(),
		cacheStorage,
		genesisInformation,
		rewardsInfo,
		spnChainID,
		lastBlockHeight,
		consumerUnbondingTime,
	)
}
