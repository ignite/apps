package formula

import (
	"math/big"

	"cosmossdk.io/math"
)

const (
	// Quadratic represents a quadratic airdrop type
	Quadratic Type = "quadratic"
)

type (
	// Value defines a struct for the formula type
	Value struct {
		Type   Type  `json:"type" yaml:"type"`
		Value  int64 `json:"value" yaml:"value"`
		Ignore int64 `json:"ignore" yaml:"ignore"`
	}
	// Type defines a formula type
	Type string
)

// Calculate calculates the airdrop amount base on the formula type
// and parameters, total amount, staked amount and the balance
func (v Value) Calculate(amount, staked math.Int) math.Int {
	switch v.Type {
	case Quadratic:
		if amount.IsZero() || amount.IsNil() {
			return math.ZeroInt()
		}

		airdrop := math.NewIntFromBigInt(big.NewInt(0).Sqrt(amount.BigInt()))

		if staked.IsPositive() {
			stakedPercent := amount.Quo(staked)
			bonus := airdrop.Mul(math.NewInt(v.Value)).Quo(stakedPercent)
			airdrop = airdrop.Add(bonus)
		}

		if airdrop.LTE(math.NewInt(v.Ignore)) {
			return math.ZeroInt()
		}

		return airdrop
	}
	return math.ZeroInt()
}
