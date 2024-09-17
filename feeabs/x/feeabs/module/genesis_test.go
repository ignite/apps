package feeabs_test

import (
	"testing"

	keepertest "feeabs/testutil/keeper"
	"feeabs/testutil/nullify"
	feeabs "feeabs/x/feeabs/module"
	"feeabs/x/feeabs/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.FeeabsKeeper(t)
	feeabs.InitGenesis(ctx, k, genesisState)
	got := feeabs.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
