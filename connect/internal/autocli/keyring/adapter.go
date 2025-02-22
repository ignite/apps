package keyring

import (
	signingv1beta1 "cosmossdk.io/api/cosmos/tx/signing/v1beta1"
	"cosmossdk.io/core/address"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

// NewAutoCLIKeyring wraps the SDK keyring and make it compatible with the AutoCLI keyring interfaces.
func NewAutoCLIKeyring(kr keyring.Keyring, ac address.Codec) (Keyring, error) {
	return &autoCLIKeyringAdapter{kr, ac}, nil
}

type autoCLIKeyringAdapter struct {
	keyring.Keyring
	ac address.Codec
}

func (a *autoCLIKeyringAdapter) List() ([]string, error) {
	list, err := a.Keyring.List()
	if err != nil {
		return nil, err
	}

	names := make([]string, len(list))
	for i, key := range list {
		names[i] = key.Name
	}

	return names, nil
}

// LookupAddressByKeyName returns the address of a key stored in the keyring
func (a *autoCLIKeyringAdapter) LookupAddressByKeyName(name string) ([]byte, error) {
	record, err := a.Keyring.Key(name)
	if err != nil {
		return nil, err
	}

	addr, err := record.GetAddress()
	if err != nil {
		return nil, err
	}

	return addr, nil
}

func (a *autoCLIKeyringAdapter) GetPubKey(name string) (cryptotypes.PubKey, error) {
	record, err := a.Keyring.Key(name)
	if err != nil {
		return nil, err
	}

	return record.GetPubKey()
}

func (a *autoCLIKeyringAdapter) Sign(name string, msg []byte, signMode signingv1beta1.SignMode) ([]byte, error) {
	record, err := a.Keyring.Key(name)
	if err != nil {
		return nil, err
	}

	signBytes, _, err := a.Keyring.Sign(record.Name, msg, signing.SignMode(signMode))
	return signBytes, err
}

func (a *autoCLIKeyringAdapter) KeyType(name string) (uint, error) {
	record, err := a.Keyring.Key(name)
	if err != nil {
		return 0, err
	}

	return uint(record.GetType()), nil
}

func (a *autoCLIKeyringAdapter) KeyInfo(nameOrAddr string) (string, string, uint, error) {
	addr, err := a.ac.StringToBytes(nameOrAddr)
	if err != nil {
		// If conversion fails, it's likely a name, not an address
		record, err := a.Keyring.Key(nameOrAddr)
		if err != nil {
			return "", "", 0, err
		}
		addr, err = record.GetAddress()
		if err != nil {
			return "", "", 0, err
		}
		addrStr, err := a.ac.BytesToString(addr)
		if err != nil {
			return "", "", 0, err
		}
		return record.Name, addrStr, uint(record.GetType()), nil
	}

	// If conversion succeeds, it's an address, get the key info by address
	record, err := a.Keyring.KeyByAddress(sdk.AccAddress(addr))
	if err != nil {
		return "", "", 0, err
	}

	return record.Name, nameOrAddr, uint(record.GetType()), nil
}
