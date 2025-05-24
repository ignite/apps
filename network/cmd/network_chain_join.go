package cmd

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/cliquiz"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/gitpod"
	"github.com/ignite/cli/v28/ignite/pkg/xchisel"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/rdegges/go-ipify"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network"
	"github.com/ignite/apps/network/network/networkchain"
)

const (
	flagGentx       = "gentx"
	flagAmount      = "amount"
	flagNoAccount   = "no-account"
	flagPeerAddress = "peer-address"
)

// NewNetworkChainJoin creates a new chain join command to join
// to a network as a validator.
func NewNetworkChainJoin() *cobra.Command {
	c := &cobra.Command{
		Use:   "join [launch-id]",
		Short: "Request to join a network as a validator",
		Long: `The "join" command is used by validators to send a request to join a blockchain.
The required argument is a launch ID of a blockchain. The "join" command expects
that the validator has already setup a home directory for the blockchain and has
a gentx either by running "ignite network chain init" or initializing the data
directory manually with the chain's binary.

By default the "join" command just sends the request to join as a validator.
However, often a validator also needs to request an genesis account with a token
balance to afford self-delegation.

The following command will send a request to join blockchain with launch ID 42
as a validator and request to be added as an account with a token balance of
95000000 STAKE.

	ignite network chain join 42 --amount 95000000stake

A request to join as a validator contains a gentx file. Ignite looks for gentx
in a home directory used by "ignite network chain init" by default. To use a
different directory, use the "--home" flag or pass a gentx file directly with
the  "--gentx" flag.

To join a chain as a validator, you must provide the IP address of your node so
that other validators can connect to it. The join command will ask you for the
IP address and will attempt to automatically detect and fill in the value. If
you want to manually specify the IP address, you can use the "--peer-address"
flag:

	ignite network chain join 42 --peer-address 0.0.0.0

Since "join" broadcasts a transaction to the Ignite blockchain, you will need an
account on the Ignite blockchain. During the testnet phase, however, Ignite
automatically requests tokens from a faucet.
`,
		Args: cobra.ExactArgs(1),
		RunE: networkChainJoinHandler,
	}

	c.Flags().String(flagGentx, "", "path to a gentx json file")
	c.Flags().String(flagAmount, "", "amount of coins for account request (ignored if coordinator has fixed the account balances or if --no-acount flag is set)")
	c.Flags().String(flagPeerAddress, "", "peer's address")
	c.Flags().Bool(flagNoAccount, false, "prevent sending a request for a genesis account")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().AddFlagSet(flagSetCheckDependencies())

	return c
}

func networkChainJoinHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	var (
		joinOptions  []network.JoinOption
		gentxPath, _ = cmd.Flags().GetString(flagGentx)
		amount, _    = cmd.Flags().GetString(flagAmount)
		noAccount, _ = cmd.Flags().GetBool(flagNoAccount)
	)

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID.
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	// if there is no custom gentx, we need to detect the public address.
	if gentxPath == "" {
		// get the peer public address for the validator.
		publicAddr, err := askPublicAddress(cmd, session)
		if err != nil {
			return err
		}

		joinOptions = append(joinOptions, network.WithPublicAddress(publicAddr))
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

	var networkOptions []networkchain.Option

	if flagGetCheckDependencies(cmd) {
		networkOptions = append(networkOptions, networkchain.CheckDependencies())
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch), networkOptions...)
	if err != nil {
		return err
	}

	// use the default gentx path from chain home if not provided
	if gentxPath == "" {
		gentxPath, err = c.DefaultGentxPath()
		if err != nil {
			return err
		}
	} else {
		// if a custom gentx is provided, we initialize the chain home in order to check accounts
		if err := c.Init(cmd.Context(), cacheStorage); err != nil {
			return err
		}
	}

	// genesis account request
	if !noAccount {
		switch {
		case c.IsAccountBalanceFixed():
			// fixed account balance
			joinOptions = append(joinOptions, network.WithAccountRequest(c.AccountBalance()))
		case amount != "":
			// account balance set by user
			amountCoins, err := sdk.ParseCoinsNormalized(amount)
			if err != nil {
				return errors.Wrap(err, "error parsing amount")
			}
			joinOptions = append(joinOptions, network.WithAccountRequest(amountCoins))
		default:
			// fixed balance and no amount entered by the user, we ask if they want to skip account request
			if !getYes(cmd) {
				question := fmt.Sprintf(
					"You haven't set the --%s flag and therefore an account request won't be submitted. Do you confirm",
					flagAmount,
				)
				if err := session.AskConfirm(question); err != nil {
					if errors.Is(err, promptui.ErrAbort) {
						return nil
					}

					return err
				}
			}

			session.Printf("%s %s\n", icons.Info, "Account request won't be submitted")
		}
	}

	// create requests to join as a validator
	joinRequests, err := n.GetJoinRequestContents(cmd.Context(), c, launchID, gentxPath, joinOptions...)
	if err != nil {
		return err
	}

	// simulate the join requests
	if err := verifyRequestsFromRequestContents(
		cmd.Context(),
		cacheStorage,
		nb,
		launchID,
		joinRequests...,
	); err != nil {
		return err
	}

	// send join requests
	return n.SendRequests(cmd.Context(), launchID, joinRequests)
}

// askPublicAddress prepare questions to interactively ask for a publicAddress
// when peer isn't provided and not running through chisel proxy.
func askPublicAddress(cmd *cobra.Command, session *cliui.Session) (publicAddress string, err error) {
	ctx := cmd.Context()

	if gitpod.IsOnGitpod() {
		publicAddress, err = gitpod.URLForPort(ctx, xchisel.DefaultServerPort)
		if err != nil {
			return "", errors.Wrap(err, "cannot read public Gitpod address of the node")
		}
		return publicAddress, nil
	}

	peerAddress, _ := cmd.Flags().GetString(flagPeerAddress)

	// The `--peer-address` flag is required when "--yes" is present
	if getYes(cmd) && peerAddress == "" {
		return "", errors.New("a peer address is required")
	}

	// Don't prompt for an address when it is available as a flag value
	if peerAddress != "" {
		return peerAddress, nil
	}

	// Try to guess the current peer address. This address is used
	// as default when prompting user for the right peer address.
	if ip, err := ipify.GetIp(); err == nil {
		peerAddress = fmt.Sprintf("%s:26656", ip)
	}

	options := []cliquiz.Option{cliquiz.Required()}
	if peerAddress != "" {
		options = append(options, cliquiz.DefaultAnswer(peerAddress))
	}

	questions := []cliquiz.Question{cliquiz.NewQuestion(
		"Peer's address",
		&publicAddress,
		options...,
	)}
	return publicAddress, session.Ask(questions...)
}
