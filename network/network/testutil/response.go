package testutil

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"google.golang.org/protobuf/runtime/protoiface"
)

// NewResponse creates cosmosclient.Response object from proto struct
// for using as a return result for a cosmosclient mock.
func NewResponse(data protoiface.MessageV1) cosmosclient.Response {
	marshaler := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	anyEncoded, _ := codectypes.NewAnyWithValue(data)

	txData := &sdk.TxMsgData{MsgResponses: []*codectypes.Any{anyEncoded}}

	encodedTxData, _ := marshaler.Marshal(txData)
	resp := cosmosclient.Response{
		Codec: marshaler,
		TxResponse: &sdk.TxResponse{
			Data: hex.EncodeToString(encodedTxData),
		},
	}
	return resp
}
