package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_relativePath(t *testing.T) {
	tests := []struct {
		name    string
		appPath string
		want    string
		err     error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := relativePath(tt.appPath)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
