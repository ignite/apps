package scaffolder

import (
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_assertSupportedCosmosSDKVersion(t *testing.T) {
	v100, err := cosmosver.Parse("v1.0.0")
	require.NoError(t, err)
	v0501, err := cosmosver.Parse("v0.50.1")
	require.NoError(t, err)
	v0450, err := cosmosver.Parse("v0.45.0")
	require.NoError(t, err)
	v0391, err := cosmosver.Parse("v0.39.1")
	require.NoError(t, err)

	tests := []struct {
		name string
		v    cosmosver.Version
		err  error
	}{
		{
			name: "Supported Cosmos SDK version (equal)",
			v:    v0501,
		},
		{
			name: "Supported Cosmos SDK version (greater than)",
			v:    v100,
		},
		{
			name: "Unsupported Cosmos SDK version",
			v:    v0450,
			err:  errors.Errorf(errOldCosmosSDKVersionStr, v0450),
		},
		{
			name: "Unsupported Cosmos SDK version",
			v:    v0391,
			err:  errors.Errorf(errOldCosmosSDKVersionStr, v0391),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assertSupportedCosmosSDKVersion(tt.v)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
