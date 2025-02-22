package keyring

import (
	"io"

	"github.com/spf13/pflag"

	"github.com/ignite/apps/connect/internal/flags"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"

	signingv1beta1 "cosmossdk.io/api/cosmos/tx/signing/v1beta1"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
)

// ContextKey is the key used to store the keyring in the context.
// The keyring must be wrapped using the KeyringImpl.
var ContextKey keyringContextKey

type keyringContextKey struct{}

var _ Keyring = &KeyringImpl{}

type KeyringImpl struct {
	k Keyring
}

// NewIgniteKeyring creates a new keyring instance based on command-line flags.
func NewIgniteKeyring(
	flagSet *pflag.FlagSet,
	input io.Reader,
) (*KeyringImpl, error) {
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

	igniteKeyring, err := keyring.NewAutoCLIKeyring(ca.Keyring)
	if err != nil {
		return nil, err
	}

	return &KeyringImpl{k: igniteKeyring}, nil
}

// GetPubKey implements Keyring.
func (k *KeyringImpl) GetPubKey(name string) (types.PubKey, error) {
	return k.k.GetPubKey(name)
}

// List implements Keyring.
func (k *KeyringImpl) List() ([]string, error) {
	return k.k.List()
}

// LookupAddressByKeyName implements Keyring.
func (k *KeyringImpl) LookupAddressByKeyName(name string) ([]byte, error) {
	return k.k.LookupAddressByKeyName(name)
}

// Sign implements Keyring.
func (k *KeyringImpl) Sign(name string, msg []byte, signMode signingv1beta1.SignMode) ([]byte, error) {
	return k.k.Sign(name, msg, signMode)
}
