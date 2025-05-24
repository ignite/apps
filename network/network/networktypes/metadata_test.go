package networktypes_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/apps/network/network/networktypes"
)

func TestMetadata_IsCurrentVersion(t *testing.T) {
	tests := []struct {
		name string
		m    networktypes.Metadata
		want bool
	}{
		{
			name: "current version",
			m: networktypes.Metadata{
				Cli: networktypes.Cli{
					Version: networktypes.Version,
				},
			},
			want: true,
		},
		{
			name: "not current version",
			m: networktypes.Metadata{
				Cli: networktypes.Cli{
					Version: "0",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.m.IsCurrentVersion())
		})
	}
}
