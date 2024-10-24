package network

import (
	"context"
	"testing"

	launchtypes "github.com/ignite/network/x/launch/types"
	"github.com/stretchr/testify/require"

	"github.com/ignite/apps/network/network/networktypes"
	"github.com/ignite/apps/network/network/testutil"
)

func TestSendRequest(t *testing.T) {
	t.Run("successfully send request", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			module         = "module"
			param          = "param"
			value          = []byte("value")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				launchtypes.NewMsgSendRequest(
					addr,
					testutil.LaunchID,
					launchtypes.NewParamChange(
						testutil.LaunchID,
						module,
						param,
						value,
					),
				),
			).
			Return(testutil.NewResponse(&launchtypes.MsgSendRequestResponse{
				RequestId:    0,
				AutoApproved: false,
			}), nil).
			Once()

		sendRequestError := network.SendRequest(context.Background(), testutil.LaunchID, launchtypes.NewParamChange(
			testutil.LaunchID,
			module,
			param,
			value,
		))
		require.NoError(t, sendRequestError)
		suite.AssertAllMocks(t)
	})
}
