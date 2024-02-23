package config

import (
	"os"
	"strings"
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAddWasm(t *testing.T) {
	type args struct {
		configPath string
		options    []Option
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Add wasm parameters to the config file",
			args: args{
				configPath: "testdata/config_without_wasm.toml",
				options: []Option{
					WithSmartQueryGasLimit(77),
					WithMemoryCacheSize(888),
					WithSimulationGasLimit(9999),
				},
			},
		},
		{
			name: "Invalid config file path",
			args: args{
				configPath: "nonexistent_directory/nonexistent_config.toml",
				options: []Option{
					WithSmartQueryGasLimit(77),
					WithMemoryCacheSize(888),
					WithSimulationGasLimit(9999),
				},
			},
			err: errors.New("open nonexistent_directory/nonexistent_config.toml: no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, _ := os.ReadFile(tt.args.configPath)
			err := AddWasm(tt.args.configPath, tt.args.options...)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.True(t, hasWasm(tt.args.configPath))

			withWasm, err := os.ReadFile("testdata/config_with_wasm.toml")
			require.NoError(t, err)
			noWasm, err := os.ReadFile(tt.args.configPath)
			require.NoError(t, err)
			require.Equal(t, strings.TrimSpace(string(withWasm)), strings.TrimSpace(string(noWasm)))

			require.NoError(t, os.WriteFile(tt.args.configPath, content, 0o644))
			require.False(t, hasWasm(tt.args.configPath))
		})
	}
}

func Test_hasWasm(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		want       bool
	}{
		{
			name:       "Config file with wasm section",
			configPath: "testdata/config_with_wasm.toml",
			want:       true,
		},
		{
			name:       "Config file without wasm section",
			configPath: "testdata/config_without_wasm.toml",
			want:       false,
		},
		{
			name:       "Non-existent config file",
			configPath: "testdata/nonexistent_config.toml",
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasWasm(tt.configPath)
			require.Equal(t, tt.want, got)
		})
	}
}
