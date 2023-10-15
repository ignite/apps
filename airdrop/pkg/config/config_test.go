package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/ignite/apps/airdrop/pkg/formula"
)

func TestParseConfig(t *testing.T) {
	sampleConfig := Config{
		AirdropToken: "ufoo",
		DustWallet:   1,
		Snapshots: []Snapshot{
			{
				Type:  "staking",
				Denom: "uatom",
				Formula: formula.Value{
					Type:  "quadratic",
					Value: 2,
				},
				Excluded: []string{"cosmos1aqn8ynvr3jmq67879qulzrwhchq5dtrvh6h4er"},
			},
			{
				Type:  "liquidity",
				Denom: "uatom",
				Formula: formula.Value{
					Type:  "quadratic",
					Value: 10,
				},
				Excluded: []string{"cosmos1aqn8ynvr3jmq67879qulzrwhchq5dtrvh6h4er"},
			},
		},
	}
	yamlData, err := yaml.Marshal(&sampleConfig)
	require.NoError(t, err)
	sampleConfigPath := filepath.Join(t.TempDir(), "config.yml")
	err = os.WriteFile(sampleConfigPath, yamlData, 0o644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filename string
		want     Config
		err      error
	}{
		{
			name:     "valid config file",
			filename: sampleConfigPath,
			want:     sampleConfig,
		},
		{
			name:     "valid config file",
			filename: "invalid_file_path",
			want:     sampleConfig,
			err:      errors.New("open invalid_file_path: no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfig(tt.filename)
			if tt.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want.AirdropToken, got.AirdropToken)
			require.Equal(t, tt.want.DustWallet, got.DustWallet)
			for i, wantSnapshot := range tt.want.Snapshots {
				require.Equal(t, wantSnapshot.Denom, got.Snapshots[i].Denom)
				require.Equal(t, wantSnapshot.Formula, got.Snapshots[i].Formula)
				require.Equal(t, wantSnapshot.Type, got.Snapshots[i].Type)
				require.EqualValues(t, wantSnapshot.Excluded, got.Snapshots[i].Excluded)

			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name string
		c    Config
		err  error
	}{
		{
			name: "nil snapshots",
			c: Config{
				AirdropToken: "uatom",
				DustWallet:   0,
				Snapshots:    nil,
			},
			err: ErrInvalidConfig,
		},
		{
			name: "empty snapshots",
			c: Config{
				AirdropToken: "uatom",
				DustWallet:   0,
				Snapshots:    []Snapshot{},
			},
			err: ErrInvalidConfig,
		},
		{
			name: "empty airdrop token",
			c: Config{
				AirdropToken: "",
				DustWallet:   0,
				Snapshots: []Snapshot{
					{
						Type:  "staking",
						Denom: "uatom",
						Formula: formula.Value{
							Type:   formula.Quadratic,
							Value:  2,
							Ignore: 1,
						},
						Excluded: nil,
					},
				},
			},
			err: ErrInvalidConfig,
		},
		{
			name: "empty snapshot type",
			c: Config{
				AirdropToken: "uatom",
				DustWallet:   0,
				Snapshots: []Snapshot{
					{
						Type:  "",
						Denom: "uatom",
						Formula: formula.Value{
							Type:   formula.Quadratic,
							Value:  2,
							Ignore: 1,
						},
						Excluded: nil,
					},
				},
			},
			err: ErrInvalidSnapshotConfig,
		},
		{
			name: "empty snapshot denom",
			c: Config{
				AirdropToken: "uatom",
				DustWallet:   0,
				Snapshots: []Snapshot{
					{
						Type:  "staking",
						Denom: "",
						Formula: formula.Value{
							Type:   formula.Quadratic,
							Value:  2,
							Ignore: 1,
						},
						Excluded: nil,
					},
				},
			},
			err: ErrInvalidSnapshotConfig,
		},
		{
			name: "empty snapshot formula type",
			c: Config{
				AirdropToken: "uatom",
				DustWallet:   0,
				Snapshots: []Snapshot{
					{
						Type:  "staking",
						Denom: "uatom",
						Formula: formula.Value{
							Type:   "",
							Value:  2,
							Ignore: 1,
						},
						Excluded: nil,
					},
				},
			},
			err: ErrInvalidSnapshotConfig,
		},
		{
			name: "zero formula value",
			c: Config{
				AirdropToken: "uatom",
				DustWallet:   0,
				Snapshots: []Snapshot{
					{
						Type:  "staking",
						Denom: "uatom",
						Formula: formula.Value{
							Type:   formula.Quadratic,
							Value:  0,
							Ignore: 1,
						},
						Excluded: nil,
					},
				},
			},
			err: ErrInvalidSnapshotConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.validate()
			if tt.err != nil {
				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), tt.err.Error()))
				return
			}
			require.NoError(t, err)
		})
	}
}
