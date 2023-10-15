package client

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
)

// Client represents the chain query client
type Client struct {
	cosmos       cosmosclient.Client
	stakingQuery stakingtypes.QueryClient
	bankQuery    banktypes.QueryClient
}

// New creates a new client for chain query API.
func New(cosmos cosmosclient.Client) Client {
	return Client{
		cosmos:       cosmos,
		stakingQuery: stakingtypes.NewQueryClient(cosmos.Context()),
		bankQuery:    banktypes.NewQueryClient(cosmos.Context()),
	}
}

func (c Client) Validators(
	ctx context.Context,
	pagination *query.PageRequest,
) (stakingtypes.Validators, error) {
	req := &stakingtypes.QueryValidatorsRequest{
		Pagination: pagination,
	}
	resp, err := c.stakingQuery.Validators(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error requesting validators: %s", err)
	}
	return resp.Validators, nil
}

func (c Client) Balances(
	ctx context.Context,
	pagination *query.PageRequest,
	denom string,
) ([]*banktypes.DenomOwner, error) {
	req := &banktypes.QueryDenomOwnersRequest{
		Denom:      denom,
		Pagination: pagination,
	}
	resp, err := c.bankQuery.DenomOwners(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error requesting denom balances: %s", err)
	}
	return resp.DenomOwners, nil
}
