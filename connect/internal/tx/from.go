package tx

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"

	sdkkeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetFromAddress gets the from address from the cobra command.
func GetFromAddress(cmd *cobra.Command) (string, error) {
	from, err := cmd.Flags().GetString(flags.FlagFrom)
	if err != nil {
		return "", err
	}

	txCtx, err := GetContext(cmd.Context())
	if err != nil {
		return "", err
	}

	fromAddr, keyType, err := GetFromFields(txCtx, from)
	if err != nil {
		return "", fmt.Errorf("failed to convert address field to address: %w", err)
	}

	if keyType == sdkkeyring.TypeLedger {
		// Set sign-mode flag to legacy Amino JSON when using a Ledger key.
		_ = cmd.Flags().Set(flags.FlagSignMode, flags.SignModeLegacyAminoJSON)
	}

	return fromAddr.String(), nil
}

// GetFromFields returns a from account address and keyring type, given either an address or key name.
func GetFromFields(ctx Context, from string) (sdk.AccAddress, sdkkeyring.KeyType, error) {
	if from == "" {
		from = ctx.Keyring.DefaultKey()
		if from == "" {
			return nil, 0, fmt.Errorf("no key name or address provided")
		}
	}

	addr, err := ctx.AddressCodec.StringToBytes(from)
	var k *sdkkeyring.Record

	sdkKeyring, ok := ctx.Keyring.Impl().(sdkkeyring.Keyring)
	if !ok {
		return nil, 0, fmt.Errorf("keyring does not implement sdkkeyring.Keyring")
	}

	if err == nil {
		k, err = sdkKeyring.KeyByAddress(sdk.AccAddress(addr))
		if err != nil {
			return nil, 0, err
		}
	} else {
		k, err = sdkKeyring.Key(from)
		if err != nil {
			return nil, 0, err
		}
	}

	addr, err = k.GetAddress()
	if err != nil {
		return nil, 0, err
	}

	return addr, k.GetType(), nil
}
