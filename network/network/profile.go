package network

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ignite/cli/v28/ignite/pkg/events"
	profiletypes "github.com/ignite/network/x/profile/types"
	projecttypes "github.com/ignite/network/x/project/types"

	"github.com/ignite/apps/network/network/networktypes"
)

// CoordinatorIDByAddress returns the CoordinatorByAddress from SPN.
func (n Network) CoordinatorIDByAddress(ctx context.Context, address string) (uint64, error) {
	n.ev.Send("Fetching coordinator by address", events.ProgressStart())
	resCoordByAddr, err := n.profileQuery.GetCoordinatorByAddress(ctx,
		&profiletypes.QueryGetCoordinatorByAddressRequest{
			Address: address,
		},
	)

	if isNotFoundErr(err) {
		return 0, ErrObjectNotFound
	} else if err != nil {
		return 0, err
	}
	return resCoordByAddr.Coordinator.CoordinatorId, nil
}

// SetCoordinatorDescription set the description of a coordinator
// or creates the coordinator if it doesn't exist yet for the sender address.
func (n Network) SetCoordinatorDescription(ctx context.Context, description profiletypes.CoordinatorDescription) error {
	n.ev.Send("Setting coordinator description", events.ProgressStart())

	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	// check if coordinator exists
	_, err = n.CoordinatorIDByAddress(ctx, addr)
	if errors.Is(err, ErrObjectNotFound) {
		// create a new coordinator
		msgCreateCoordinator := profiletypes.NewMsgCreateCoordinator(
			addr,
			description.Identity,
			description.Website,
			description.Details,
		)
		res, err := n.cosmos.BroadcastTx(ctx, n.account, msgCreateCoordinator)
		if err != nil {
			return err
		}
		var requestRes profiletypes.MsgCreateCoordinatorResponse
		return res.Decode(&requestRes)
	} else if err == nil {
		// update the description for the coordinator
		msgUpdateCoordinatorDescription := profiletypes.NewMsgUpdateCoordinatorDescription(
			addr,
			description.Identity,
			description.Website,
			description.Details,
		)
		res, err := n.cosmos.BroadcastTx(ctx, n.account, msgUpdateCoordinatorDescription)
		if err != nil {
			return err
		}
		var requestRes profiletypes.MsgUpdateCoordinatorDescriptionResponse
		return res.Decode(&requestRes)
	}
	return err
}

// Coordinator returns the Coordinator by address from SPN.
func (n Network) Coordinator(ctx context.Context, address string) (networktypes.Coordinator, error) {
	n.ev.Send("Fetching coordinator details", events.ProgressStart())
	coordinatorID, err := n.CoordinatorIDByAddress(ctx, address)
	if err != nil {
		return networktypes.Coordinator{}, err
	}
	resCoord, err := n.profileQuery.GetCoordinator(ctx,
		&profiletypes.QueryGetCoordinatorRequest{
			CoordinatorId: coordinatorID,
		},
	)
	if isNotFoundErr(err) {
		return networktypes.Coordinator{}, ErrObjectNotFound
	} else if err != nil {
		return networktypes.Coordinator{}, err
	}
	return networktypes.ToCoordinator(resCoord.Coordinator), nil
}

// SetValidatorDescription set a validator profile.
func (n Network) SetValidatorDescription(ctx context.Context, validator profiletypes.Validator) error {
	n.ev.Send("Setting validator description", events.ProgressStart())

	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	message := profiletypes.NewMsgUpdateValidatorDescription(
		addr,
		validator.Description.Identity,
		validator.Description.Moniker,
		validator.Description.Website,
		validator.Description.SecurityContact,
		validator.Description.Details,
	)

	res, err := n.cosmos.BroadcastTx(ctx, n.account, message)
	if err != nil {
		return err
	}

	var requestRes profiletypes.MsgUpdateValidatorDescriptionResponse
	return res.Decode(&requestRes)
}

// Validator returns the Validator by address from SPN.
func (n Network) Validator(ctx context.Context, address string) (networktypes.Validator, error) {
	n.ev.Send("Fetching validator description", events.ProgressStart())
	res, err := n.profileQuery.GetValidator(ctx, &profiletypes.QueryGetValidatorRequest{
		Address: address,
	})
	if isNotFoundErr(err) {
		return networktypes.Validator{}, ErrObjectNotFound
	} else if err != nil {
		return networktypes.Validator{}, err
	}
	return networktypes.ToValidator(res.Validator), nil
}

// Balances returns the all balances by address from SPN.
func (n Network) Balances(ctx context.Context, address string) (sdk.Coins, error) {
	n.ev.Send("Fetching address balances", events.ProgressStart())
	res, err := banktypes.NewQueryClient(n.cosmos.Context()).AllBalances(ctx,
		&banktypes.QueryAllBalancesRequest{
			Address: address,
		},
	)
	if isNotFoundErr(err) {
		return sdk.Coins{}, ErrObjectNotFound
	} else if err != nil {
		return sdk.Coins{}, err
	}
	return res.Balances, nil
}

// Profile returns the address profile info.
func (n Network) Profile(ctx context.Context, projectID uint64) (networktypes.Profile, error) {
	address, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return networktypes.Profile{}, err
	}

	// fetch vouchers held by the account
	coins, err := n.Balances(ctx, address)
	if err != nil {
		return networktypes.Profile{}, err
	}
	vouchers := sdk.NewCoins()
	for _, coin := range coins {
		// parse the coin to filter all non-voucher coins from the balance
		_, err := projecttypes.VoucherProject(coin.Denom)
		if err == nil {
			vouchers = append(vouchers, coin)
		}
	}
	vouchers = vouchers.Sort()

	var shares projecttypes.Shares

	// if a project ID is specified, fetches the shares of the project
	if projectID > 0 {
		acc, err := n.MainnetAccount(ctx, projectID, address)
		if err != nil && !errors.Is(err, ErrObjectNotFound) {
			return networktypes.Profile{}, err
		}
		shares = acc.Shares
	}

	var p networktypes.ProfileAcc
	p, err = n.Validator(ctx, address)
	if errors.Is(err, ErrObjectNotFound) {
		p, err = n.Coordinator(ctx, address)
		if errors.Is(err, ErrObjectNotFound) {
			p = networktypes.Coordinator{Address: address}
		} else if err != nil {
			return networktypes.Profile{}, err
		}
	} else if err != nil {
		return networktypes.Profile{}, err
	}
	return p.ToProfile(projectID, vouchers, shares), nil
}
