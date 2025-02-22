package tx

import (
	// "fmt"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	// "github.com/cosmos/cosmos-sdk/crypto/keyring"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

/* TODO */

// GetFromAddress gets the from address from the cobra command.
// It checks the from flags, as well as the potential default ignite account (the first one)
func GetFromAddress(cmd *cobra.Command) (string, error) {
	return "", nil
}

// func setFrom() {
// 	if clientCtx.From == "" || flagSet.Changed(flags.FlagFrom) {
// 		from, _ := flagSet.GetString(flags.FlagFrom)
// 		fromAddr, fromName, keyType, err := GetFromFields(clientCtx, clientCtx.Keyring, from)
// 		if err != nil {
// 			return clientCtx, fmt.Errorf("failed to convert address field to address: %w", err)
// 		}

// 		clientCtx = clientCtx.WithFrom(from).WithFromAddress(fromAddr).WithFromName(fromName)

// 		if keyType == keyring.TypeLedger && clientCtx.SignModeStr == flags.SignModeTextual {
// 			if !slicsliceses.Contains(clientCtx.TxConfig.SignModeHandler().SupportedModes(), signingv1beta1.SignMode_SIGN_MODE_TEXTUAL) {
// 				return clientCtx, fmt.Errorf("SIGN_MODE_TEXTUAL is not available")
// 			}
// 		}

// 		// If the `from` signer account is a ledger key, we need to use
// 		// SIGN_MODE_AMINO_JSON, because ledger doesn't support proto yet.
// 		// ref: https://github.com/cosmos/cosmos-sdk/issues/8109
// 		if keyType == keyring.TypeLedger &&
// 			clientCtx.SignModeStr != flags.SignModeLegacyAminoJSON &&
// 			clientCtx.SignModeStr != flags.SignModeTextual &&
// 			!clientCtx.LedgerHasProtobuf {
// 			fmt.Println("Default sign-mode 'direct' not supported by Ledger, using sign-mode 'amino-json'.")
// 			clientCtx = clientCtx.WithSignModeStr(flags.SignModeLegacyAminoJSON)
// 		}
// 	}
// }

// // GetFromFields returns a from account address, account name and keyring type, given either an address or key name.
// // If clientCtx.Simulate is true the keystore is not accessed and a valid address must be provided
// // If clientCtx.GenerateOnly is true the keystore is only accessed if a key name is provided
// // If from is empty, the default key if specified in the context will be used
// func GetFromFields(clientCtx Context, kr keyring.Keyring, from string) (sdk.AccAddress, string, keyring.KeyType, error) {
// 	if from == "" && clientCtx.KeyringDefaultKeyName != "" {
// 		from = clientCtx.KeyringDefaultKeyName
// 		_ = clientCtx.PrintString(fmt.Sprintf("No key name or address provided; using the default key: %s\n", clientCtx.KeyringDefaultKeyName))
// 	}

// 	if from == "" {
// 		return nil, "", 0, nil
// 	}

// 	addr, err := clientCtx.AddressCodec.StringToBytes(from)
// 	switch {
// 	case clientCtx.Simulate:
// 		if err != nil {
// 			return nil, "", 0, fmt.Errorf("a valid address must be provided in simulation mode: %w", err)
// 		}

// 		return addr, "", 0, nil

// 	case clientCtx.GenerateOnly:
// 		if err == nil {
// 			return addr, "", 0, nil
// 		}
// 	}

// 	var k *keyring.Record
// 	if err == nil {
// 		k, err = kr.KeyByAddress(addr)
// 		if err != nil {
// 			return nil, "", 0, err
// 		}
// 	} else {
// 		k, err = kr.Key(from)
// 		if err != nil {
// 			return nil, "", 0, err
// 		}
// 	}

// 	addr, err = k.GetAddress()
// 	if err != nil {
// 		return nil, "", 0, err
// 	}

// 	return addr, k.Name, k.GetType(), nil
// }
