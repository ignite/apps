package network

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/pkg/xurl"
	launchtypes "github.com/ignite/network/x/launch/types"

	"github.com/ignite/apps/network/network/address"
	"github.com/ignite/apps/network/network/gentx"
	"github.com/ignite/apps/network/network/networkchain"
	"github.com/ignite/apps/network/network/networktypes"
)

type joinOptions struct {
	accountAmount sdk.Coins
	publicAddress string
}

type JoinOption func(*joinOptions)

// WithAccountRequest allows to join the chain by requesting a genesis account with the specified amount of tokens.
func WithAccountRequest(amount sdk.Coins) JoinOption {
	return func(o *joinOptions) {
		o.accountAmount = amount
	}
}

// WithPublicAddress allows to specify a peer public address for the node.
func WithPublicAddress(addr string) JoinOption {
	return func(o *joinOptions) {
		o.publicAddress = addr
	}
}

// GetJoinRequestContents returns the request contents to join a chain as a validator.
func (n Network) GetJoinRequestContents(
	ctx context.Context,
	c Chain,
	launchID uint64,
	gentxPath string,
	options ...JoinOption,
) (reqs []launchtypes.RequestContent, err error) {
	o := joinOptions{}
	for _, apply := range options {
		apply(&o)
	}

	var (
		nodeID string
		peer   launchtypes.Peer
	)

	// parse the gentx content
	gentxInfo, gentx, err := gentx.GentxFromPath(gentxPath)
	if err != nil {
		return reqs, err
	}

	// get the peer address
	if o.publicAddress != "" {
		if nodeID, err = c.NodeID(ctx); err != nil {
			return reqs, err
		}

		if xurl.IsHTTP(o.publicAddress) {
			peer = launchtypes.NewPeerTunnel(nodeID, networkchain.HTTPTunnelChisel, o.publicAddress)
		} else {
			peer = launchtypes.NewPeerConn(nodeID, o.publicAddress)
		}
	} else {
		// if the peer address is not specified, we parse it from the gentx memo
		if peer, err = ParsePeerAddress(gentxInfo.Memo); err != nil {
			return reqs, err
		}
	}

	// change the chain address prefix to spn
	accountAddress, err := address.ChangeValidatorAddressPrefix(gentxInfo.ValidatorAddress, networktypes.SPN)
	if err != nil {
		return reqs, err
	}

	if !o.accountAmount.IsZero() {
		reqs = append(reqs, launchtypes.NewGenesisAccount(
			launchID,
			accountAddress,
			o.accountAmount,
		))
	}

	reqs = append(reqs, launchtypes.NewGenesisValidator(
		launchID,
		accountAddress,
		gentx,
		gentxInfo.PubKey,
		gentxInfo.SelfDelegation,
		peer,
	))

	return reqs, nil
}
