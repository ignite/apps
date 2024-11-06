package cmd

import (
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/xurl"
	"github.com/ignite/network/pkg/chainid"
	launchtypes "github.com/ignite/network/x/launch/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network"
	"github.com/ignite/apps/network/network/networkchain"
)

const (
	flagTag            = "tag"
	flagBranch         = "branch"
	flagHash           = "hash"
	flagGenesisURL     = "genesis-url"
	flagGenesisConfig  = "genesis-config"
	flagProject        = "project"
	flagShares         = "shares"
	flagNoCheck        = "no-check"
	flagChainID        = "chain-id"
	flagMainnet        = "mainnet"
	flagAccountBalance = "account-balance"
	flagRewardCoins    = "reward.coins"
	flagRewardHeight   = "reward.height"
)

// NewNetworkChainPublish returns a new command to publish a new chain to start a new network.
func NewNetworkChainPublish() *cobra.Command {
	c := &cobra.Command{
		Use:   "publish [source-url]",
		Short: "Publish a new chain to start a new network",
		Long: `To begin the process of launching a blockchain with Ignite, a coordinator needs
to publish the information about a blockchain. The only required bit of
information is the URL of the source code of the blockchain.

The following command publishes the information about an example blockchain:

	ignite network chain publish github.com/ignite/example

This command fetches the source code of the blockchain, compiles the binary,
verifies that a blockchain can be started with the binary, and publishes the
information about the blockchain to Ignite. Currently, only public repositories
are supported. The command returns an integer number that acts as an identifier
of the chain on Ignite.

By publishing a blockchain on Ignite you become the "coordinator" of this
blockchain. A coordinator is an account that has the authority to approve and
reject validator requests, set parameters of the blockchain and trigger the
launch of the chain.

The default Git branch is used when publishing a chain. If you want to use a
specific branch, tag or a commit hash, use "--branch", "--tag", or "--hash"
flags respectively.

The repository name is used as the default chain ID. Ignite does not ensure that
chain IDs are unique, but they have to have a valid format: [string]-[integer].
To set a custom chain ID use the "--chain-id" flag.

	ignite network chain publish github.com/ignite/example --chain-id foo-1

Once the chain is published users can request accounts with coin balances to be
added to the chain's genesis. By default, users are free to request any number
of tokens. If you want all users requesting tokens to get the same amount, use
the "--account-balance" flag with a list of coins.

	ignite network chain publish github.com/ignite/example --account-balance 2000foocoin
`,
		Args: cobra.ExactArgs(1),
		RunE: networkChainPublishHandler,
	}

	flagSetClearCache(c)
	c.Flags().String(flagBranch, "", "Git branch to use for the repo")
	c.Flags().String(flagTag, "", "Git tag to use for the repo")
	c.Flags().String(flagHash, "", "Git hash to use for the repo")
	c.Flags().String(flagGenesisURL, "", "URL to a custom Genesis")
	c.Flags().String(flagGenesisConfig, "", "name of an Ignite config file in the repo for custom Genesis")
	c.Flags().String(flagChainID, "", "chain ID to use for this network")
	c.Flags().Int64(flagProject, -1, "project ID to use for this network")
	c.Flags().Bool(flagNoCheck, false, "skip verifying chain's integrity")
	c.Flags().String(flagMetadata, "", "add chain metadata")
	c.Flags().String(flagProjectTotalSupply, "", "add a total of the mainnet of a project")
	c.Flags().String(flagShares, "", "add shares for the project")
	c.Flags().Bool(flagMainnet, false, "initialize a mainnet project")
	c.Flags().String(flagAccountBalance, "", "balance for each approved genesis account for the chain")
	c.Flags().String(flagRewardCoins, "", "reward coins")
	c.Flags().Int64(flagRewardHeight, 0, "last reward height")
	c.Flags().String(flagAmount, "", "amount of coins for account request")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().AddFlagSet(flagSetCheckDependencies())

	return c
}

func networkChainPublishHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	var (
		tag, _                   = cmd.Flags().GetString(flagTag)
		branch, _                = cmd.Flags().GetString(flagBranch)
		hash, _                  = cmd.Flags().GetString(flagHash)
		genesisURL, _            = cmd.Flags().GetString(flagGenesisURL)
		genesisConfig, _         = cmd.Flags().GetString(flagGenesisConfig)
		chainID, _               = cmd.Flags().GetString(flagChainID)
		project, _               = cmd.Flags().GetInt64(flagProject)
		noCheck, _               = cmd.Flags().GetBool(flagNoCheck)
		metadata, _              = cmd.Flags().GetString(flagMetadata)
		projectTotalSupplyStr, _ = cmd.Flags().GetString(flagProjectTotalSupply)
		sharesStr, _             = cmd.Flags().GetString(flagShares)
		isMainnet, _             = cmd.Flags().GetBool(flagMainnet)
		accountBalance, _        = cmd.Flags().GetString(flagAccountBalance)
		rewardCoinsStr, _        = cmd.Flags().GetString(flagRewardCoins)
		rewardDuration, _        = cmd.Flags().GetInt64(flagRewardHeight)
		amount, _                = cmd.Flags().GetString(flagAmount)
	)

	// parse the amount.
	amountCoins, err := sdk.ParseCoinsNormalized(amount)
	if err != nil {
		return errors.Wrap(err, "error parsing amount")
	}

	accountBalanceCoins, err := sdk.ParseCoinsNormalized(accountBalance)
	if err != nil {
		return errors.Wrap(err, "error parsing account balance")
	}

	source, err := xurl.MightHTTPS(args[0])
	if err != nil {
		return fmt.Errorf("invalid source url format: %w", err)
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	if project > -1 && projectTotalSupplyStr != "" {
		return fmt.Errorf("%s and %s flags cannot be set together", flagProject, flagProjectTotalSupply)
	}
	if isMainnet {
		if project < 0 && projectTotalSupplyStr == "" {
			return fmt.Errorf(
				"%s flag requires one of the %s or %s flags to be set",
				flagMainnet,
				flagProject,
				flagProjectTotalSupply,
			)
		}
		if chainID == "" {
			return fmt.Errorf("%s flag requires the %s flag", flagMainnet, flagChainID)
		}
	}

	if chainID != "" {
		chainName, _, err := chainid.ParseGenesisChainID(chainID)
		if err != nil {
			return errors.Wrapf(err, "invalid chain id: %s", chainID)
		}
		if err := chainid.CheckChainName(chainName); err != nil {
			return errors.Wrapf(err, "invalid chain id name: %s", chainName)
		}
	}

	totalSupply, err := sdk.ParseCoinsNormalized(projectTotalSupplyStr)
	if err != nil {
		return err
	}

	rewardCoins, err := sdk.ParseCoinsNormalized(rewardCoinsStr)
	if err != nil {
		return err
	}

	if (!rewardCoins.Empty() && rewardDuration == 0) ||
		(rewardCoins.Empty() && rewardDuration > 0) {
		return fmt.Errorf("%s and %s flags must be provided together", flagRewardCoins, flagRewardHeight)
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// use source from chosen target.
	var sourceOption networkchain.SourceOption

	switch {
	case tag != "":
		sourceOption = networkchain.SourceRemoteTag(source, tag)
	case branch != "":
		sourceOption = networkchain.SourceRemoteBranch(source, branch)
	case hash != "":
		sourceOption = networkchain.SourceRemoteHash(source, hash)
	default:
		sourceOption = networkchain.SourceRemote(source)
	}

	var initOptions []networkchain.Option

	// cannot use both genesisURL and genesisConfig
	if genesisURL != "" && genesisConfig != "" {
		return errors.New("cannot use both genesis-url and genesis-config for initial genesis." +
			"Please use only one of the options.")
	}

	// use custom genesis from url if given.
	if genesisURL != "" {
		initOptions = append(initOptions, networkchain.WithGenesisFromURL(genesisURL))
	}

	// use custom genesis config if given
	if genesisConfig != "" {
		initOptions = append(initOptions, networkchain.WithGenesisFromConfig(genesisConfig))
	}

	// init in a temp dir.
	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(homeDir)

	initOptions = append(initOptions, networkchain.WithHome(homeDir))

	// prepare publish options
	publishOptions := []network.PublishOption{network.WithMetadata(metadata)}

	switch {
	case genesisURL != "":
		publishOptions = append(publishOptions, network.WithCustomGenesisURL(genesisURL))
	case genesisConfig != "":
		publishOptions = append(publishOptions, network.WithCustomGenesisConfig(genesisConfig))

	}

	if project > -1 {
		publishOptions = append(publishOptions, network.WithProject(project))
	} else if projectTotalSupplyStr != "" {
		totalSupply, err := sdk.ParseCoinsNormalized(projectTotalSupplyStr)
		if err != nil {
			return err
		}
		if !totalSupply.Empty() {
			publishOptions = append(publishOptions, network.WithTotalSupply(totalSupply))
		}
	}

	// use custom chain id if given.
	if chainID != "" {
		publishOptions = append(publishOptions, network.WithChainID(chainID))
	}

	if !accountBalanceCoins.IsZero() {
		publishOptions = append(publishOptions, network.WithAccountBalance(accountBalanceCoins))
	}

	if isMainnet {
		publishOptions = append(publishOptions, network.Mainnet())
	}

	if !totalSupply.Empty() {
		publishOptions = append(publishOptions, network.WithTotalSupply(totalSupply))
	}

	if sharesStr != "" {
		sharePercentages, err := network.ParseSharePercents(sharesStr)
		if err != nil {
			return err
		}

		publishOptions = append(publishOptions, network.WithPercentageShares(sharePercentages))
	}

	// TODO: Issue an error or warning when this flag is used with "no-check"?
	//       The "check-dependencies" flag is ignored when the "no-check" one is present.
	if flagGetCheckDependencies(cmd) {
		initOptions = append(initOptions, networkchain.CheckDependencies())
	}

	// init the chain.
	c, err := nb.Chain(sourceOption, initOptions...)
	if err != nil {
		return err
	}

	if !noCheck {
		if err := c.Init(cmd.Context(), cacheStorage); err != nil {
			// initialize the chain for checking.
			return fmt.Errorf("blockchain init failed: %w", err)
		}
	}

	session.StartSpinner("Publishing...")

	n, err := nb.Network()
	if err != nil {
		return err
	}

	launchID, projectID, err := n.Publish(cmd.Context(), c, publishOptions...)
	if err != nil {
		return err
	}

	if !rewardCoins.IsZero() && rewardDuration > 0 {
		if err := n.SetReward(cmd.Context(), launchID, rewardDuration, rewardCoins); err != nil {
			return err
		}
	}

	if !amountCoins.IsZero() {
		// create a request to add an account to the genesis
		addr, err := n.AccountAddress()
		if err != nil {
			return err
		}
		addAccountRequest := launchtypes.NewGenesisAccount(
			launchID,
			addr,
			amountCoins,
		)

		// simulate the add account request
		if err := verifyRequestsFromRequestContents(
			cmd.Context(),
			cacheStorage,
			nb,
			launchID,
			addAccountRequest,
		); err != nil {
			return err
		}

		// send the request
		if err := n.SendRequest(cmd.Context(), launchID, addAccountRequest); err != nil {
			return err
		}
	}

	session.Printf("%s Network published \n", icons.OK)
	if isMainnet {
		session.Printf("%s Mainnet ID: %d \n", icons.Bullet, launchID)
	} else {
		session.Printf("%s Launch ID: %d \n", icons.Bullet, launchID)
	}
	if projectID > -1 {
		session.Printf("%s Project ID: %d \n", icons.Bullet, projectID)
	}

	return nil
}
