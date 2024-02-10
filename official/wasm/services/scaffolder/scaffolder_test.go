package scaffolder

import (
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosver"
	"github.com/stretchr/testify/require"
)

func Test_assertSupportedCosmosSDKVersion(t *testing.T) {
	tests := []struct {
		name string
		v    cosmosver.Version
		err  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assertSupportedCosmosSDKVersion(tt.v)
			require.Equal(t, tt.err, err)
		})
	}
}
