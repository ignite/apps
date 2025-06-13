package scaffolder

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
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
		name        string
		wasmVersion semver.Version
		sdkVersion  cosmosver.Version
		want        bool
		err         error
	}{
		{
			name:        "Supported Cosmos SDK version (equal)",
			wasmVersion: DefaultWasmVersion,
			sdkVersion:  v0501,
			err:         errors.Errorf(errNotCompatibleVersionStr, DefaultWasmVersion, v0501),
		},
		{
			name:        "Supported legacy Wasm version",
			wasmVersion: LegacyWasmVersion,
			sdkVersion:  v100,
			want:        true,
		},
		{
			name:        "Unsupported Wasm version",
			wasmVersion: LegacyWasmVersion,
			sdkVersion:  v0450,
			err:         errors.Errorf(errOldCosmosSDKVersionStr, v0450),
		},
		{
			name:        "Unsupported Wasm version",
			wasmVersion: semver.MustParse("0.60.0"),
			sdkVersion:  v0450,
			err:         errors.Errorf(errOldCosmosSDKVersionStr, v0450),
		},
		{
			name:        "Unsupported Cosmos SDK version",
			wasmVersion: DefaultWasmVersion,
			sdkVersion:  v0391,
			err:         errors.Errorf(errOldCosmosSDKVersionStr, v0391),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assertVersions(tt.wasmVersion, tt.sdkVersion)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
