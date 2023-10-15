package snapshot

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ignite/apps/airdrop/pkg/encode"
)

// Generate produce the airdrop snapshot based on chain
// state balances and stakes values
func Generate(genState map[string]json.RawMessage) (Snapshot, error) {
	var (
		marshaller   = encode.Codec()
		snapshotAccs = make(Accounts)
	)

	var bankGenesis banktypes.GenesisState
	if len(genState[banktypes.ModuleName]) > 0 {
		err := marshaller.UnmarshalJSON(genState[banktypes.ModuleName], &bankGenesis)
		if err != nil {
			return Snapshot{}, err
		}
	}
	for _, balance := range bankGenesis.Balances {
		var (
			address = balance.Address
			acc     = snapshotAccs.getAccount(address)
		)
		acc.Balance = balance.Coins
		snapshotAccs[address] = acc
	}

	var stakingGenesis stakingtypes.GenesisState
	if len(genState[stakingtypes.ModuleName]) > 0 {
		err := marshaller.UnmarshalJSON(genState[stakingtypes.ModuleName], &stakingGenesis)
		if err != nil {
			return Snapshot{}, err
		}
	}
	for _, unbonding := range stakingGenesis.UnbondingDelegations {
		var (
			address        = unbonding.DelegatorAddress
			acc            = snapshotAccs.getAccount(address)
			unbondingStake = sdk.NewInt(0)
		)

		for _, entry := range unbonding.Entries {
			unbondingStake = unbondingStake.Add(entry.Balance)
		}

		acc.UnbondingStake = acc.UnbondingStake.Add(unbondingStake)
		snapshotAccs[address] = acc
	}

	// Make a map from validator operator address to the v036 validator type
	validators := make(map[string]stakingtypes.Validator)
	for _, validator := range stakingGenesis.Validators {
		validators[validator.OperatorAddress] = validator
	}

	for _, delegation := range stakingGenesis.Delegations {
		var (
			address = delegation.DelegatorAddress
			acc     = snapshotAccs.getAccount(address)
			val     = validators[delegation.ValidatorAddress]
			staked  = delegation.Shares.MulInt(val.Tokens).Quo(val.DelegatorShares).RoundInt()
		)
		acc.Staked = acc.Staked.Add(staked)
		snapshotAccs[address] = acc
	}

	return Snapshot{
		NumberAccounts: uint64(len(snapshotAccs)),
		BondDenom:      stakingGenesis.Params.BondDenom,
		Accounts:       snapshotAccs,
	}, nil
}
