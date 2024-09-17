package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "feeabs/testutil/keeper"
	"feeabs/x/feeabs/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.FeeabsKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
