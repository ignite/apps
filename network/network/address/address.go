package address

import (
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

// ChangeValidatorAddressPrefix returns the address with another prefix from the validator address.
func ChangeValidatorAddressPrefix(addr, newPrefix string) (string, error) {
	return ChangeAddressPrefix(addr, "cosmosvaloper", newPrefix)
}

// ChangeAddressPrefix returns the address with another prefix.
func ChangeAddressPrefix(addr, prefix, newPrefix string) (string, error) {
	if newPrefix == "" {
		return "", errors.New("empty new prefix")
	}
	if prefix == "" {
		return "", errors.New("empty prefix")
	}
	cdc := address.NewBech32Codec(prefix)
	bAddr, err := cdc.StringToBytes(addr)
	if err != nil {
		return "", err
	}
	return bech32.ConvertAndEncode(newPrefix, bAddr)
}
