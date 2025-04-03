package internal

import (
	"github.com/spf13/pflag"

	"cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/ignite/apps/connect/internal/flags"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
)

// NewKeyring creates a new keyring instance based on command-line flags.
func NewKeyring(
	flagSet *pflag.FlagSet,
	addressCodec address.Codec,
	bech32Prefix string,
) (keyring.Keyring, error) {
	keyringBackend, err := flagSet.GetString(flags.FlagKeyringBackend)
	if err != nil {
		return nil, err
	} else if keyringBackend == "" {
		keyringBackend = keyring.BackendTest
	}

	ca, err := cosmosaccount.New(
		cosmosaccount.WithBech32Prefix(bech32Prefix),
		cosmosaccount.WithHome(cosmosaccount.KeyringHome),
		cosmosaccount.WithKeyringBackend(cosmosaccount.KeyringBackend(keyringBackend)),
	)
	if err != nil {
		return nil, err
	}

	return ca.Keyring, nil
}
