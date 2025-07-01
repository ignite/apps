package scaffolder

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func Test_assertSupportedCosmosSDKVersion(t *testing.T) {
	v0391, err := cosmosver.Parse("v0.39.1")
	require.NoError(t, err)
	v0449, err := cosmosver.Parse("v0.44.9")
	require.NoError(t, err)
	v0450, err := cosmosver.Parse("v0.45.0")
	require.NoError(t, err)
	v0500, err := cosmosver.Parse("v0.50.0")
	require.NoError(t, err)
	v0501, err := cosmosver.Parse("v0.50.1")
	require.NoError(t, err)
	v0510, err := cosmosver.Parse("0.51.0")
	require.NoError(t, err)
	v100, err := cosmosver.Parse("v1.0.0")
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
			name: "Unsupported Cosmos SDK version",
			v:    v0450,
			err:  errors.Errorf(errOldCosmosSDKVersionStr, v0450),
		},
		{
			name: "Unsupported Cosmos SDK version",
			v:    v0391,
			err:  errors.Errorf(errOldCosmosSDKVersionStr, v0391),
		},

		{
			name: "Supported Cosmos SDK version (equal to v0.50.1)",
			v:    v0501,
		},
		{
			name: "Supported Cosmos SDK version (greater than v1.0.0)",
			v:    v100,
			err:  errors.Errorf(errNewCosmosSDKVersionStr, v100, v0510),
		},
		{
			name: "Unsupported Cosmos SDK version (less than v0.45.0)",
			v:    v0391,
			err:  errors.Errorf(errOldCosmosSDKVersionStr, v0391),
		},
		{
			name: "Unsupported Cosmos SDK version (less than v0.50.1)",
			v:    v0450,
			err:  errors.Errorf(errOldCosmosSDKVersionStr, v0450),
		},
		{
			name: "Edge case: exact boundary of supported version v0.50.1",
			v:    v0501,
		},
		{
			name: "Edge case: exact boundary of unsupported version v0.39.1",
			v:    v0391,
			err:  errors.Errorf(errOldCosmosSDKVersionStr, v0391),
		},
		{
			name: "Lower boundary case: supported version v0.50.0 (if supported)",
			v:    v0500,
			err:  errors.Errorf(errOldCosmosSDKVersionStr, v0500),
		},
		{
			name: "Lower boundary case: unsupported version v0.44.9 (if unsupported)",
			v:    v0449,
			err:  errors.Errorf(errOldCosmosSDKVersionStr, v0449),
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
