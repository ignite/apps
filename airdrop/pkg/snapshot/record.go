package snapshot

import (
	"cosmossdk.io/math"
	claimtypes "github.com/ignite/modules/x/claim/types"

	"github.com/ignite/apps/airdrop/pkg/formula"
)

type (
	// Record provide a record with all airdrop balances
	Record map[string]claimtypes.ClaimRecord

	// ConfigType represents a Filter type
	ConfigType string

	// Records represents an array of Filter's
	Records []Record
)

const (
	// Staking filter type staking
	Staking = "staking"
	// Liquidity filter type liquidity
	Liquidity = "liquidity"
)

// ClaimRecords return a list of claim records
func (f Record) ClaimRecords() []claimtypes.ClaimRecord {
	result := make([]claimtypes.ClaimRecord, 0)
	for _, filter := range f {
		result = append(result, filter)
	}
	return result
}

// Sum sum all filters into one
func (f Records) Sum() Record {
	result := make(Record)
	for _, filter := range f {
		for _, amount := range filter {
			resultAmount := result.getAmount(amount.Address)
			resultAmount.Claimable = resultAmount.Claimable.Add(amount.Claimable)
			result[amount.Address] = resultAmount
		}
	}
	return result
}

// getAccount get an existing account or generate a new one
func (f Record) getAmount(address string) claimtypes.ClaimRecord {
	acc, ok := f[address]
	if ok {
		return acc
	}
	return claimtypes.ClaimRecord{
		Address:   address,
		Claimable: math.NewInt(0),
	}
}

// ApplyConfig apply the configuration to the snaplshot filtering based
// on the config type, denom and excluded address, and apply the formula to generate
// the snapshot
func (s Snapshot) ApplyConfig(
	configType ConfigType,
	denom string,
	formula formula.Value,
	excludedAddresses []string,
) Record {
	if len(excludedAddresses) > 0 {
		s.Accounts.excludeAddresses(excludedAddresses...)
	}
	s.Accounts.filterDenom(denom)

	filter := make(Record)
	for address, account := range s.Accounts {
		amount := account.balanceAmount()
		if configType == Staking {
			amount = account.balanceAmount()
		}
		claimAmount := formula.Calculate(amount, account.Staked)
		filter[address] = claimtypes.ClaimRecord{
			Address:   address,
			Claimable: claimAmount,
		}
	}
	return filter
}
