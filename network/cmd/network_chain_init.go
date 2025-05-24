package cmd

import (
	"errors"
	"fmt"

	"github.com/ignite/cli/v28/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/cliquiz"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	cosmosgenesis "github.com/ignite/cli/v28/ignite/pkg/cosmosutil/genesis"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network"
	"github.com/ignite/apps/network/network/networkchain"
)

const (
	flagValidatorAccount         = "validator-account"
	flagValidatorWebsite         = "validator-website"
	flagValidatorDetails         = "validator-details"
	flagValidatorSecurityContact = "validator-security-contact"
	flagValidatorMoniker         = "validator-moniker"
	flagValidatorIdentity        = "validator-identity"
	flagValidatorSelfDelegation  = "validator-self-delegation"
	flagValidatorGasPrice        = "validator-gas-price"

	defaultSelfDelegation          = "95000000stake"
	defaultCommissionMax           = "0.10"
	defaultCommissionMaxRate       = "0.20"
	defaultCommissionMaxChangeRate = "0.01"
)

// NewNetworkChainInit returns a new command to initialize a chain from a published chain ID.
func NewNetworkChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init [launch-id]",
		Short: "Initialize a chain from a published chain ID",
		Long: `Ignite network chain init is a command used by validators to initialize a
validator node for a blockchain from the information stored on the Ignite chain.

	ignite network chain init 42

This command fetches the information about a chain with launch ID 42. The source
code of the chain is cloned in a temporary directory, and the node's binary is
compiled from the source. The binary is then used to initialize the node. By
default, Ignite uses "~/spn/[launch-id]/" as the home directory for the blockchain.

An important part of initializing a validator node is creation of the gentx (a
transaction that adds a validator at the genesis of the chain).

The "init" command will prompt for values like self-delegation and commission.
These values will be used in the validator's gentx. You can use flags to provide
the values in non-interactive mode.

Use the "--home" flag to choose a different path for the home directory of the
blockchain:

	ignite network chain init 42 --home ~/mychain

The end result of the "init" command is a validator home directory with a
genesis validator transaction (gentx) file.`,
		Args: cobra.ExactArgs(1),
		RunE: networkChainInitHandler,
	}

	flagSetClearCache(c)
	c.Flags().String(flagValidatorAccount, cosmosaccount.DefaultAccount, "account for the chain validator")
	c.Flags().String(flagValidatorWebsite, "", "associate a website with the validator")
	c.Flags().String(flagValidatorDetails, "", "details about the validator")
	c.Flags().String(flagValidatorSecurityContact, "", "validator security contact email")
	c.Flags().String(flagValidatorMoniker, "", "custom validator moniker")
	c.Flags().String(flagValidatorIdentity, "", "validator identity signature (ex. UPort or Keybase)")
	c.Flags().String(flagValidatorSelfDelegation, "", "validator minimum self delegation")
	c.Flags().String(flagValidatorGasPrice, "", "validator gas price")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	return c
}

func networkChainInitHandler(cmd *cobra.Command, args []string) error {
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

	// check if the provided account for the validator exists.
	validatorAccount, _ := cmd.Flags().GetString(flagValidatorAccount)
	if _, err = nb.AccountRegistry.GetByName(validatorAccount); err != nil {
		return err
	}

	// if a chain has already been initialized with this launch ID, we ask for confirmation
	// before erasing the directory.
	chainHome, exist, err := networkchain.IsChainHomeExist(launchID)
	if err != nil {
		return err
	}

	if !getYes(cmd) && exist {
		question := fmt.Sprintf(
			"The chain has already been initialized under: %s. Would you like to overwrite the home directory",
			chainHome,
		)
		if err := session.AskConfirm(question); err != nil {
			if errors.Is(err, promptui.ErrAbort) {
				return nil
			}

			return err
		}
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
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

	if err := c.Init(cmd.Context(), cacheStorage); err != nil {
		return err
	}

	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}

	genesis, err := cosmosgenesis.FromPath(genesisPath)
	if err != nil {
		return err
	}
	stakeDenom, err := genesis.StakeDenom()
	if err != nil {
		return err
	}

	// ask validator information.
	defaultSelfDel := chainLaunch.AccountBalance.String()
	if defaultSelfDel == "" {
		defaultSelfDel = defaultSelfDelegation
	}
	v, err := askValidatorInfo(cmd, session, stakeDenom, defaultSelfDel)
	if err != nil {
		return err
	}
	session.StartSpinner("Generating your Gentx")

	gentxPath, err := c.InitAccount(cmd.Context(), v, validatorAccount)
	if err != nil {
		return err
	}

	return session.Printf("%s Gentx generated: %s\n", icons.Bullet, gentxPath)
}

// askValidatorInfo prompts to the user questions to query validator information.
func askValidatorInfo(
	cmd *cobra.Command,
	session *cliui.Session,
	stakeDenom,
	defaultSelfDel string,
) (chain.Validator, error) {
	var (
		account, _         = cmd.Flags().GetString(flagValidatorAccount)
		website, _         = cmd.Flags().GetString(flagValidatorWebsite)
		details, _         = cmd.Flags().GetString(flagValidatorDetails)
		securityContact, _ = cmd.Flags().GetString(flagValidatorSecurityContact)
		moniker, _         = cmd.Flags().GetString(flagValidatorMoniker)
		identity, _        = cmd.Flags().GetString(flagValidatorIdentity)
		selfDelegation, _  = cmd.Flags().GetString(flagValidatorSelfDelegation)
		gasPrice, _        = cmd.Flags().GetString(flagValidatorGasPrice)
	)
	if gasPrice == "" {
		gasPrice = "0" + stakeDenom
	}
	v := chain.Validator{
		Name:              account,
		Website:           website,
		Details:           details,
		Moniker:           moniker,
		Identity:          identity,
		SecurityContact:   securityContact,
		MinSelfDelegation: selfDelegation,
		GasPrices:         gasPrice,
	}

	questions := append([]cliquiz.Question{},
		cliquiz.NewQuestion("Staking amount",
			&v.StakingAmount,
			cliquiz.DefaultAnswer(defaultSelfDel),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Commission rate",
			&v.CommissionRate,
			cliquiz.DefaultAnswer(defaultCommissionMax),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Commission max rate",
			&v.CommissionMaxRate,
			cliquiz.DefaultAnswer(defaultCommissionMaxRate),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Commission max change rate",
			&v.CommissionMaxChangeRate,
			cliquiz.DefaultAnswer(defaultCommissionMaxChangeRate),
			cliquiz.Required(),
		),
	)
	return v, session.Ask(questions...)
}
