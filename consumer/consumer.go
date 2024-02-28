package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	pluginv1 "github.com/ignite/cli/v28/ignite/services/plugin/grpc/v1"

	cmtypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	cmprivval "github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibctmtypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	ccvconsumertypes "github.com/cosmos/interchain-security/v3/x/ccv/consumer/types"
	ccvtypes "github.com/cosmos/interchain-security/v3/x/ccv/types"
)

// writeConsumerGenesis writes the consumer module genesis in the genesis file.
func writeConsumerGenesis(chain *pluginv1.ChainInfo) error {
	var (
		providerClientState = &ibctmtypes.ClientState{
			ChainId:         "provider",
			TrustLevel:      ibctmtypes.DefaultTrustLevel,
			TrustingPeriod:  time.Hour * 64,
			UnbondingPeriod: time.Hour * 128,
			MaxClockDrift:   time.Minute * 5,
		}
		providerConsState = &ibctmtypes.ConsensusState{
			Timestamp: time.Now().Add(time.Hour * 24),
			Root: commitmenttypes.NewMerkleRoot(
				[]byte("LpGpeyQVLUo9HpdsgJr12NP2eCICspcULiWa5u9udOA="),
			),
			NextValidatorsHash: []byte("E30CE736441FB9101FADDAF7E578ABBE6DFDB67207112350A9A904D554E1F5BE"),
		}
		params = ccvtypes.NewParams(
			true,
			1000, // ignore distribution
			"",   // ignore distribution
			"",   // ignore distribution
			ccvtypes.DefaultCCVTimeoutPeriod,
			ccvtypes.DefaultTransferTimeoutPeriod,
			ccvtypes.DefaultConsumerRedistributeFrac,
			ccvtypes.DefaultHistoricalEntries,
			ccvtypes.DefaultConsumerUnbondingPeriod,
			"0", // disable soft opt-out
			[]string{},
			[]string{},
		)
	)
	// Load public key from priv_validator_key.json file
	pk, err := getPubKey(chain)
	if err != nil {
		return err
	}
	// Feed initial_val_set with this public key
	// Like for sovereign chain, provide only a single validator.
	valUpdates := cmtypes.ValidatorUpdates{
		cmtypes.UpdateValidator(pk.Bytes(), 1, pk.Type()),
	}
	// Build consumer genesis
	consumerGen := ccvtypes.NewInitialGenesisState(providerClientState, providerConsState, valUpdates, params)
	// Read genesis file
	genPath := getGenesisPath(chain)
	genState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genPath)
	if err != nil {
		return err
	}
	// Update consumer module gen state
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	bz, err := codec.MarshalJSON(consumerGen)
	if err != nil {
		return err
	}
	genState[ccvconsumertypes.ModuleName] = bz
	// Update whole genesis
	bz, err = json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return err
	}
	genDoc.AppState = bz
	// Save genesis
	return genDoc.SaveAs(genPath)
}

// isInitialized returns true if the consumer chain `chain` is initizalied.
// A consumer chain is considered initialized if its consumer genesis contains
// at least one validator in the InitialValSet field.
func isInitialized(chain *pluginv1.ChainInfo) (bool, error) {
	genPath := getGenesisPath(chain)
	genState, _, err := genutiltypes.GenesisStateFromGenFile(genPath)
	if err != nil {
		// If the genesis isn't readable, don't propagate the error, just
		// consider the chain isn't initialized.
		return false, nil
	}
	var (
		consumerGenesis   ccvtypes.GenesisState
		interfaceRegistry = codectypes.NewInterfaceRegistry()
		codec             = codec.NewProtoCodec(interfaceRegistry)
	)
	err = codec.UnmarshalJSON(genState[ccvconsumertypes.ModuleName], &consumerGenesis)
	if err != nil {
		return false, err
	}
	return len(consumerGenesis.InitialValSet) != 0, nil
}

// getPubKey returns the validator public key.
func getPubKey(chain *pluginv1.ChainInfo) (crypto.PubKey, error) {
	keyFilePath := filepath.Join(chain.Home, "config", "priv_validator_key.json")
	keyJSONBytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	var pvKey cmprivval.FilePVKey
	err = cmtjson.Unmarshal(keyJSONBytes, &pvKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading PrivValidator key from %v: %w", keyFilePath, err)
	}
	return pvKey.PubKey, nil
}

// getGenesisPath returns genesis.json path of the chain.
func getGenesisPath(chain *pluginv1.ChainInfo) string {
	return filepath.Join(chain.Home, "config", "genesis.json")
}
