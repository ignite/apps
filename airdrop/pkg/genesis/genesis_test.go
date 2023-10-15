package genesis

import (
	"testing"

	claimtypes "github.com/ignite/modules/x/claim/types"
	"github.com/stretchr/testify/require"
)

func TestGenState_AddFromClaimRecord(t *testing.T) {
	tests := []struct {
		name    string
		g       GenState
		denom   string
		records []claimtypes.ClaimRecord
		err     error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.g.AddFromClaimRecord(tt.denom, tt.records)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, tt.err, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetGenStateFromPath(t *testing.T) {
	tests := []struct {
		name    string
		genesis string
		want    GenState
		err     error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGenStateFromPath(tt.genesis)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, tt.err, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
