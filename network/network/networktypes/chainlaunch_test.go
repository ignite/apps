package networktypes_test

import (
	"testing"
	"time"

	launchtypes "github.com/ignite/network/x/launch/types"
	"github.com/stretchr/testify/require"

	"github.com/ignite/apps/network/network/networktypes"
)

func TestToChainLaunch(t *testing.T) {
	tests := []struct {
		name     string
		fetched  launchtypes.Chain
		expected networktypes.ChainLaunch
	}{
		{
			name: "chain with default genesis",
			fetched: launchtypes.Chain{
				LaunchId:       1,
				GenesisChainId: "foo-1",
				SourceUrl:      "foo.com",
				SourceHash:     "0xaaa",
				HasProject:     true,
				ProjectId:      1,
				InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
			},
			expected: networktypes.ChainLaunch{
				ID:              1,
				ChainID:         "foo-1",
				SourceURL:       "foo.com",
				SourceHash:      "0xaaa",
				GenesisURL:      "",
				GenesisHash:     "",
				LaunchTriggered: false,
				ProjectID:       1,
				Network:         "testnet",
			},
		},
		{
			name: "launched chain with custom genesis url and no project",
			fetched: launchtypes.Chain{
				LaunchId:        1,
				GenesisChainId:  "bar-1",
				SourceUrl:       "bar.com",
				SourceHash:      "0xbbb",
				LaunchTriggered: true,
				LaunchTime:      time.Unix(100, 100).UTC(),
				InitialGenesis: launchtypes.NewGenesisURL(
					"genesisfoo.com",
					"0xccc",
				),
			},
			expected: networktypes.ChainLaunch{
				ID:              1,
				ChainID:         "bar-1",
				SourceURL:       "bar.com",
				SourceHash:      "0xbbb",
				GenesisURL:      "genesisfoo.com",
				GenesisHash:     "0xccc",
				LaunchTriggered: true,
				LaunchTime:      time.Unix(100, 100).UTC(),
				ProjectID:       0,
				Network:         "testnet",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.expected, networktypes.ToChainLaunch(tt.fetched))
		})
	}
}
