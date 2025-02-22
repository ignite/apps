package tx

import (
	"testing"

	"github.com/stretchr/testify/require"

	bankv1beta1 "cosmossdk.io/api/cosmos/bank/v1beta1"
	base "cosmossdk.io/api/cosmos/base/v1beta1"
	apisigning "cosmossdk.io/api/cosmos/tx/signing/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

func getWrappedTx(t *testing.T) *wrappedTx {
	t.Helper()

	f, err := NewFactory(keybase, cdc, mockAccountRetriever{}, txConf, ac, mockClientConn{}, TxParameters{})
	require.NoError(t, err)

	pk := secp256k1.GenPrivKey().PubKey()
	addr, _ := ac.BytesToString(pk.Address())

	f.tx.msgs = []sdk.Msg{&bankv1beta1.MsgSend{
		FromAddress: addr,
		ToAddress:   addr,
		Amount:      []*base.Coin{},
	}}
	require.NoError(t, err)

	err = f.setFeePayer(addr)
	require.NoError(t, err)

	f.tx.fees = []*base.Coin{{
		Denom:  "cosmos",
		Amount: "1000",
	}}

	err = f.setSignatures([]Signature{{
		PubKey: pk,
		Data: &SingleSignatureData{
			SignMode:  apisigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: 0,
	}}...)
	require.NoError(t, err)
	wTx, err := f.getTx()
	require.NoError(t, err)
	return wTx
}

func Test_txEncoder_txDecoder(t *testing.T) {
	wTx := getWrappedTx(t)

	encodedTx, err := encodeTx(wTx)
	require.NoError(t, err)
	require.NotNil(t, encodedTx)

	isDeterministic, err := encodeTx(wTx)
	require.NoError(t, err)
	require.NotNil(t, encodedTx)
	require.Equal(t, encodedTx, isDeterministic)

	f := decodeTx(cdc, decoder)
	decodedTx, err := f(encodedTx)
	require.NoError(t, err)
	require.NotNil(t, decodedTx)

	dTx, ok := decodedTx.(*wrappedTx)
	require.True(t, ok)
	require.Equal(t, wTx.TxRaw, dTx.TxRaw)
	require.Equal(t, wTx.Tx.AuthInfo.String(), dTx.Tx.AuthInfo.String())
	require.Equal(t, wTx.Tx.Body.String(), dTx.Tx.Body.String())
	require.Equal(t, wTx.Tx.Signatures, dTx.Tx.Signatures)
}

func Test_txJsonEncoder_txJsonDecoder(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "json encode and decode tx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wTx := getWrappedTx(t)

			encodedTx, err := encodeJsonTx(wTx)
			require.NoError(t, err)
			require.NotNil(t, encodedTx)

			f := decodeJsonTx(cdc, decoder)
			decodedTx, err := f(encodedTx)
			require.NoError(t, err)
			require.NotNil(t, decodedTx)

			dTx, ok := decodedTx.(*wrappedTx)
			require.True(t, ok)
			require.Equal(t, wTx.TxRaw, dTx.TxRaw)
			require.Equal(t, wTx.Tx.AuthInfo.String(), dTx.Tx.AuthInfo.String())
			require.Equal(t, wTx.Tx.Body.String(), dTx.Tx.Body.String())
			require.Equal(t, wTx.Tx.Signatures, dTx.Tx.Signatures)
		})
	}
}
