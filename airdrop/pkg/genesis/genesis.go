package genesis

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/math"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	claimtypes "github.com/ignite/modules/x/claim/types"
	"github.com/pkg/errors"
)

type (
	// GenState defines the initial conditions for a tendermint blockchain, in particular its validator set.
	GenState struct {
		GenesisTime     time.Time                  `json:"genesis_time"`
		ChainID         string                     `json:"chain_id"`
		InitialHeight   int64                      `json:"initial_height"`
		ConsensusParams *tmproto.ConsensusParams   `json:"consensus_params,omitempty"`
		Validators      []tmtypes.GenesisValidator `json:"validators,omitempty"`
		AppHash         tmbytes.HexBytes           `json:"app_hash"`
		AppState        AppState                   `json:"app_state,omitempty"`
	}

	// AppState defines a genesis app state.
	AppState map[string]json.RawMessage
)

// GetGenStateFromPath returns a JSON genesis state message from inputted path.
func GetGenStateFromPath(genesisFilePath string) (genState GenState, err error) {
	genesisFile, err := os.Open(filepath.Clean(genesisFilePath))
	if err != nil {
		return genState, err
	}
	defer genesisFile.Close()

	byteValue, err := io.ReadAll(genesisFile)
	if err != nil {
		return genState, err
	}

	return genState, tmjson.Unmarshal(byteValue, &genState)
}

// AddFromClaimRecord add a claim record to the genesis state.
func (g GenState) AddFromClaimRecord(denom string, claimRecords []claimtypes.ClaimRecord) error {
	claimGenesis := claimtypes.GenesisState{
		ClaimRecords:  claimRecords,
		AirdropSupply: sdk.NewCoin(denom, math.ZeroInt()),
	}
	for _, claimRecord := range claimRecords {
		claimGenesis.AirdropSupply.Add(sdk.NewCoin(denom, claimRecord.Claimable))
	}
	if len(g.AppState[claimtypes.ModuleName]) > 0 {
		return errors.New("claim record state already exist into the genesis")
	}

	claimBytes, err := json.Marshal(claimGenesis)
	if err != nil {
		return err
	}
	g.AppState[claimtypes.ModuleName] = claimBytes
	return nil
}
