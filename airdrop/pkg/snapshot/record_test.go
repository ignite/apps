package snapshot

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	claimtypes "github.com/ignite/modules/x/claim/types"
	"github.com/stretchr/testify/require"

	"github.com/ignite/apps/airdrop/pkg/formula"
)

func TestAccounts_getAmount(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	sampleAmounts := Record{
		accAddr1: {
			Address:   accAddr1,
			Claimable: math.NewInt(10),
		},
		accAddr2: {
			Address:   accAddr2,
			Claimable: math.NewInt(1000),
		},
	}
	tests := []struct {
		name    string
		a       Record
		address string
		want    claimtypes.ClaimRecord
	}{
		{
			name:    "already exist address 1",
			a:       sampleAmounts,
			address: accAddr1,
			want:    sampleAmounts[accAddr1],
		},
		{
			name:    "already exist address 2",
			a:       sampleAmounts,
			address: accAddr2,
			want:    sampleAmounts[accAddr2],
		},
		{
			name:    "not exist address",
			a:       sampleAmounts,
			address: accAddr3,
			want: claimtypes.ClaimRecord{
				Address:   accAddr3,
				Claimable: math.ZeroInt(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.getAmount(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFilters_Sum(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	tests := []struct {
		name    string
		filters Records
		want    Record
	}{
		{
			name: "2 filters",
			filters: Records{
				{
					accAddr1: {
						Address:   accAddr1,
						Claimable: math.NewInt(120),
					},
					accAddr2: {
						Address:   accAddr2,
						Claimable: math.NewInt(440),
					},
				},
				{
					accAddr1: {
						Address:   accAddr1,
						Claimable: math.NewInt(224),
					},
					accAddr2: {
						Address:   accAddr2,
						Claimable: math.NewInt(233),
					},
					accAddr3: {
						Address:   accAddr3,
						Claimable: math.NewInt(233),
					},
				},
			},
			want: Record{
				accAddr1: {
					Address:   accAddr1,
					Claimable: math.NewInt(344),
				},
				accAddr2: {
					Address:   accAddr2,
					Claimable: math.NewInt(673),
				},
				accAddr3: {
					Address:   accAddr3,
					Claimable: math.NewInt(233),
				},
			},
		},
		{
			name: "3 filters",
			filters: Records{
				{
					accAddr1: {
						Address:   accAddr1,
						Claimable: math.NewInt(30),
					},
				},
				{
					accAddr1: {
						Address:   accAddr1,
						Claimable: math.NewInt(120),
					},
					accAddr2: {
						Address:   accAddr2,
						Claimable: math.NewInt(220),
					},
				},
				{
					accAddr1: {
						Address:   accAddr1,
						Claimable: math.NewInt(224),
					},
					accAddr2: {
						Address:   accAddr2,
						Claimable: math.NewInt(220),
					},
					accAddr3: {
						Address:   accAddr3,
						Claimable: math.NewInt(233),
					},
				},
			},
			want: Record{
				accAddr1: {
					Address:   accAddr1,
					Claimable: math.NewInt(374),
				},
				accAddr2: {
					Address:   accAddr2,
					Claimable: math.NewInt(440),
				},
				accAddr3: {
					Address:   accAddr3,
					Claimable: math.NewInt(233),
				},
			},
		},
		{
			name: "2 filters different addresses",
			filters: Records{
				{
					accAddr1: {
						Address:   accAddr1,
						Claimable: math.NewInt(120),
					},
				},
				{
					accAddr2: {
						Address:   accAddr2,
						Claimable: math.NewInt(220),
					},
				},
			},
			want: Record{
				accAddr1: {
					Address:   accAddr1,
					Claimable: math.NewInt(120),
				},
				accAddr2: {
					Address:   accAddr2,
					Claimable: math.NewInt(220),
				},
			},
		},
		{
			name: "2 filters same addresses",
			filters: Records{
				{
					accAddr1: {
						Address:   accAddr1,
						Claimable: math.NewInt(321),
					},
				},
				{
					accAddr1: {
						Address:   accAddr1,
						Claimable: math.NewInt(123),
					},
				},
			},
			want: Record{
				accAddr1: {
					Address:   accAddr1,
					Claimable: math.NewInt(444),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filters.Sum()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFilter_ClaimRecords(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	tests := []struct {
		name string
		f    Record
		want []claimtypes.ClaimRecord
	}{
		{
			name: "list of one claim record",
			f: Record{
				accAddr1: claimtypes.ClaimRecord{
					Address: accAddr1,
				},
			},
			want: []claimtypes.ClaimRecord{
				{Address: accAddr1},
			},
		},
		{
			name: "list of one claim record",
			f: Record{
				accAddr1: claimtypes.ClaimRecord{
					Address: accAddr1,
				},
				accAddr2: claimtypes.ClaimRecord{
					Address: accAddr2,
				},
			},
			want: []claimtypes.ClaimRecord{
				{Address: accAddr1},
				{Address: accAddr2},
			},
		},
		{
			name: "list of one claim record",
			f: Record{
				accAddr1: claimtypes.ClaimRecord{
					Address: accAddr1,
				},
				accAddr2: claimtypes.ClaimRecord{
					Address: accAddr2,
				},
				accAddr3: claimtypes.ClaimRecord{
					Address: accAddr3,
				},
			},
			want: []claimtypes.ClaimRecord{
				{Address: accAddr1},
				{Address: accAddr2},
				{Address: accAddr3},
			},
		},
		{
			name: "empty list",
			f:    Record{},
			want: []claimtypes.ClaimRecord{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.ClaimRecords()
			require.EqualValues(t, tt.want, got)
		})
	}
}

func TestSnapshot_ApplyConfig(t *testing.T) {
	type args struct {
		filterType        ConfigType
		denom             string
		formula           formula.Value
		excludedAddresses []string
	}
	tests := []struct {
		name     string
		snapshot Snapshot
		args     args
		want     Record
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.snapshot.ApplyConfig(tt.args.filterType, tt.args.denom, tt.args.formula, tt.args.excludedAddresses)
			require.Equal(t, tt.want, got)
		})
	}
}
