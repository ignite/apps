package keyring

import (
	"io"

	"github.com/spf13/pflag"

	"github.com/ignite/apps/connect/internal/flags"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"

	"cosmossdk.io/core/address"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

// NewKeyring creates a new keyring instance based on command-line flags.
func NewKeyring(
	flagSet *pflag.FlagSet,
	input io.Reader,
	addressCodec address.Codec,
) (Keyring, error) {
	keyringBackend, err := flagSet.GetString(flags.FlagKeyringBackend)
	if err != nil {
		return nil, err
	} else if keyringBackend == "" {
		keyringBackend = keyring.BackendTest
	}

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(cosmosaccount.KeyringBackend(keyringBackend)),
		cosmosaccount.WithKeyringServiceName("ignitekeyring"),
	)
	if err != nil {
		return nil, err
	}

	igniteKeyring, err := NewAutoCLIKeyring(ca.Keyring, addressCodec)
	if err != nil {
		return nil, err
	}

	return igniteKeyring, nil
}
