package snapshot

import (
	"encoding/json"
	"os"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// Snapshot provide a snapshot with all genesis Accounts
	Snapshot struct {
		NumberAccounts uint64   `json:"num_accounts" yaml:"num_accounts"`
		BondDenom      string   `json:"bond_denom" yaml:"bond_denom"`
		Accounts       Accounts `json:"accounts" yaml:"accounts"`
	}

	// Account provide fields of snapshot per account
	// It is the simplified struct we are presenting
	// in this 'balances from state export' snapshot for people.
	Account struct {
		Address        string    `json:"address" yaml:"address"`
		Staked         math.Int  `json:"staked" yaml:"staked"`
		UnbondingStake math.Int  `json:"unbonding_stake" yaml:"unbonding_stake"`
		Balance        sdk.Coins `json:"unstake" yaml:"unstake"`
	}

	// Accounts represents a map of snapshot Accounts
	Accounts map[string]Account
)

// newAccount returns a new account.
func newAccount(address string) Account {
	return Account{
		Address:        address,
		Staked:         math.ZeroInt(),
		UnbondingStake: math.ZeroInt(),
		Balance:        sdk.NewCoins(),
	}
}

// ParseSnapshot expects to find and parse a snapshot file.
func ParseSnapshot(filename string) (c Snapshot, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return c, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return c, err
	}
	return c, nil
}

// totalStake returns a sum of stake and unbounding stake
func (a Account) totalStake() math.Int {
	if a.Staked.IsNil() {
		return a.UnbondingStake
	}
	if a.UnbondingStake.IsNil() {
		return a.Staked
	}
	return a.Staked.Add(a.UnbondingStake)
}

// balanceAmount returns a sum of all denom balances
func (a Account) balanceAmount() math.Int {
	amount := math.NewInt(0)
	if !a.Staked.IsNil() {
		amount = amount.Add(a.Staked)
	}
	if !a.UnbondingStake.IsNil() {
		amount = amount.Add(a.UnbondingStake)
	}
	for _, balance := range a.Balance {
		amount = amount.Add(balance.Amount)
	}
	return amount
}

// getAccount get an existing account or generate a new one
func (a Accounts) getAccount(address string) Account {
	acc, ok := a[address]
	if ok {
		return acc
	}
	return newAccount(address)
}

// excludeAddress exclude an address from the accounts
func (a Accounts) excludeAddress(address string) {
	for accAddress := range a {
		if accAddress == address {
			delete(a, accAddress)
		}
	}
}

// excludeAddresses exclude an address list from the accounts
func (a Accounts) excludeAddresses(addresses ...string) {
	for _, address := range addresses {
		a.excludeAddress(address)
	}
}

// filterDenom filter balance by denom
func (a Accounts) filterDenom(denom string) {
	for address, account := range a {
		found, liquidBalance := account.Balance.Find(denom)
		if found {
			account.Balance = sdk.NewCoins(liquidBalance)
		} else {
			account.Balance = sdk.NewCoins()
		}
		a[address] = account
	}
}
