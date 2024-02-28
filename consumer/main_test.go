package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v3/x/ccv/consumer/types"
	ccvtypes "github.com/cosmos/interchain-security/v3/x/ccv/types"

	"github.com/ignite/cli/v28/ignite/services/plugin"
	v1 "github.com/ignite/cli/v28/ignite/services/plugin/grpc/v1"
	"github.com/ignite/cli/v28/ignite/services/plugin/mocks"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name           string
		arg            string
		expectedError  string
		expectedOutput string
		setup          func(*testing.T, string)
	}{
		{
			name:          "fail: wrong arg",
			arg:           "wrong",
			expectedError: "invalid argument \"wrong\"",
		},
		{
			name:          "fail: writeGenesis w/o priv_validator_key.json",
			arg:           "writeGenesis",
			expectedError: "open .*/config/priv_validator_key.json: no such file or directory",
		},
		{
			name: "fail: writeGenesis w/o genesis.json",
			arg:  "writeGenesis",
			setup: func(t *testing.T, path string) {
				// Add priv_validator_key.json to path
				bz, err := os.ReadFile("testdata/config/priv_validator_key.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "priv_validator_key.json"), bz, 0o777)
				require.NoError(t, err)
			},
			expectedError: ".*/config/genesis.json does not exist, run `init` first",
		},
		{
			name: "ok: writeGenesis",
			arg:  "writeGenesis",
			setup: func(t *testing.T, path string) {
				// Add priv_validator_key.json to path
				bz, err := os.ReadFile("testdata/config/priv_validator_key.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "priv_validator_key.json"), bz, 0o777)
				require.NoError(t, err)

				// Add genesis.json to path
				bz, err = os.ReadFile("testdata/config/genesis.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "genesis.json"), bz, 0o777)
				require.NoError(t, err)
			},
		},
		{
			name:           "ok: isInitialized returns false",
			arg:            "isInitialized",
			expectedOutput: "false",
		},
		{
			name: "ok: isInitialized returns true",
			arg:  "isInitialized",
			setup: func(t *testing.T, path string) {
				// isInitialized returns true if there's a consumer genesis with an
				// InitialValSet length != 0
				// Add priv_validator_key.json to path
				bz, err := os.ReadFile("testdata/config/priv_validator_key.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "priv_validator_key.json"), bz, 0o777)
				require.NoError(t, err)

				// Add genesis.json to path
				bz, err = os.ReadFile("testdata/config/genesis.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "genesis.json"), bz, 0o777)
				require.NoError(t, err)

				// Call writeGenesis to create the genesis
				err = writeConsumerGenesis(&v1.ChainInfo{Home: path})
				require.NoError(t, err)
			},
			expectedOutput: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			homePath := t.TempDir()
			err := os.MkdirAll(filepath.Join(homePath, "config"), 0o777)
			require.NoError(t, err)
			clientAPI := mocks.NewPluginClientAPI(t)
			clientAPI.EXPECT().GetChainInfo(ctx).Return(&v1.ChainInfo{
				Home: homePath,
			}, nil)
			if tt.setup != nil {
				tt.setup(t, homePath)
			}
			// Capture os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err = app{}.Execute(ctx, &plugin.ExecutedCommand{
				Args: []string{tt.arg},
			}, clientAPI)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Regexp(t, tt.expectedError, err.Error())
				return
			}
			require.NoError(t, err)
			w.Close()
			out, _ := io.ReadAll(r)
			require.Equal(t, tt.expectedOutput, string(out))
			if tt.arg == "writeGenesis" {
				// Verify genesis
				genPath := filepath.Join(homePath, "config", "genesis.json")
				genState, _, err := genutiltypes.GenesisStateFromGenFile(genPath)
				require.NoError(t, err)
				bz, ok := genState[ccvconsumertypes.ModuleName]
				require.True(t, ok, "%s module not found in genesis", ccvconsumertypes.ModuleName)
				_ = bz
				interfaceRegistry := codectypes.NewInterfaceRegistry()
				codec := codec.NewProtoCodec(interfaceRegistry)
				var gen ccvtypes.GenesisState
				codec.MustUnmarshalJSON(bz, &gen)
				require.Equal(t, "provider", gen.GetProviderClientState().ChainId)
				require.NotEmpty(t, gen.InitialValSet)
			}
		})
	}
}
