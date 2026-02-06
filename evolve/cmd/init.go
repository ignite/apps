package cmd

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	cmprivval "github.com/cometbft/cometbft/privval"
	cmttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/viper"

	evconfig "github.com/evstack/ev-node/pkg/config"
	evgenesis "github.com/evstack/ev-node/pkg/genesis"

	configchain "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

const defaultValPower = 1

func InitHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	return initEVABCI(ctx, cmd, true)
}

func initEVABCI(
	ctx context.Context,
	cmd *plugin.ExecutedCommand,
	initChain bool,
) error {
	flags := plugin.Flags(cmd.Flags)

	session := cliui.New()
	defer session.End()

	appPath, err := flags.GetString(flagPath)
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return err
	}

	rc, err := chain.New(absPath, chain.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	if initChain {
		// use val power to set validator power in ignite config.yaml
		igniteConfig, err := rc.Config()
		if err != nil {
			return err
		}

		coins := sdk.NewCoin(igniteConfig.DefaultDenom, sdkmath.NewInt((defaultValPower * int64(math.Pow10(6)))))
		igniteConfig.Validators[0].Bonded = coins.String()
		for i, account := range igniteConfig.Accounts {
			if account.Name == igniteConfig.Validators[0].Name {
				igniteConfig.Accounts[i].Coins = []string{coins.String()}
			}
		}

		if err := configchain.Save(*igniteConfig, rc.ConfigPath()); err != nil {
			return err
		}

		if err := rc.Init(ctx, chain.InitArgsAll); err != nil {
			return err
		}
	}

	home, err := rc.Home()
	if err != nil {
		return err
	}

	// modify genesis (add sequencer)
	genesisPath, err := rc.GenesisPath()
	if err != nil {
		return err
	}

	genesis, err := genutiltypes.AppGenesisFromFile(genesisPath)
	if err != nil {
		return err
	}

	pubKey, err := getPubKey(home)
	if err != nil {
		return err
	}

	genesis.Consensus.Validators = []cmttypes.GenesisValidator{
		{
			Name:    "EV-node Sequencer",
			Address: pubKey.Address(),
			PubKey:  pubKey,
			Power:   defaultValPower,
		},
	}

	if err := genesis.SaveAs(genesisPath); err != nil {
		return err
	}

	// Add DAEpochForcedInclusion field to genesis
	fieldTag, err := getJSONTag(evgenesis.Genesis{}, "DAEpochForcedInclusion")
	if err != nil {
		return fmt.Errorf("failed to get JSON tag for DAEpochForcedInclusion in Evolve genesis: %w", err)
	}

	if err := prependFieldToGenesis(genesisPath, fieldTag, "25"); err != nil {
		return fmt.Errorf("failed to add %s to genesis: %w", fieldTag, err)
	}

	// modify evolve config (add da namespace)
	evolveConfigPath := filepath.Join(home, evconfig.AppConfigDir, evconfig.ConfigName)
	evolveViper := viper.New()
	evolveViper.SetConfigFile(evolveConfigPath)
	evolveViper.ReadInConfig()

	evolveConfig, err := evconfig.LoadFromViper(evolveViper)
	if err != nil {
		return err
	}
	evolveConfig.RootDir = home

	chainID, err := rc.ID()
	if err != nil {
		return err
	}
	evolveConfig.DA.Namespace = chainID
	evolveConfig.DA.DataNamespace = fmt.Sprintf("%s-data", chainID)
	evolveConfig.DA.ForcedInclusionNamespace = fmt.Sprintf("%s-fi-txs", chainID)

	if err := evolveConfig.SaveAsYaml(); err != nil {
		return err
	}

	return session.Printf("🗃 Initialized. Checkout your evolve chain's home (data) directory: %s\n", colors.Info(home))
}

// getPubKey returns the validator public key.
func getPubKey(chainHome string) (crypto.PubKey, error) {
	keyFilePath := filepath.Join(chainHome, "config", "priv_validator_key.json")
	keyJSONBytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	var pvKey cmprivval.FilePVKey
	if err = cmtjson.Unmarshal(keyJSONBytes, &pvKey); err != nil {
		return nil, errors.Errorf("error reading PrivValidator key from %v: %s", keyFilePath, err)
	}
	return pvKey.PubKey, nil
}

// getJSONTag extracts the JSON tag value for a given field name using reflection.
func getJSONTag(v interface{}, fieldName string) (string, error) {
	t := reflect.TypeOf(v)
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return "", fmt.Errorf("field %s not found in type %s", fieldName, t.Name())
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return "", fmt.Errorf("field %s does not have a json tag", fieldName)
	}

	// Handle tags like "field_name,omitempty" by taking only the first part
	if idx := strings.Index(jsonTag, ","); idx != -1 {
		jsonTag = jsonTag[:idx]
	}

	return jsonTag, nil
}

// prependFieldToGenesis adds a field as the first field in the genesis JSON file
// without loading the entire JSON structure into memory.
func prependFieldToGenesis(genesisPath, fieldName, fieldValue string) error {
	data, err := os.ReadFile(genesisPath)
	if err != nil {
		return fmt.Errorf("failed to read genesis file: %w", err)
	}

	// Find the first opening brace
	openBraceIdx := bytes.IndexByte(data, '{')
	if openBraceIdx == -1 {
		return fmt.Errorf("invalid genesis file: no opening brace found")
	}

	// Start right after the opening brace
	insertPos := openBraceIdx + 1

	// Find where the next non-whitespace content starts
	contentStart := insertPos
	for contentStart < len(data) && (data[contentStart] == ' ' || data[contentStart] == '\n' || data[contentStart] == '\r' || data[contentStart] == '\t') {
		contentStart++
	}

	// Check if there's any content (not just closing brace)
	hasContent := contentStart < len(data) && data[contentStart] != '}'

	// Build the new field with proper formatting
	var newField string
	if hasContent {
		// There's existing content, add comma after our field
		newField = fmt.Sprintf("\n  \"%s\": %s,", fieldName, fieldValue)
	} else {
		// Empty object, no comma needed
		newField = fmt.Sprintf("\n  \"%s\": %s\n", fieldName, fieldValue)
	}

	// Construct the new file content
	var buf bytes.Buffer
	buf.Write(data[:insertPos])
	buf.WriteString(newField)
	// Write the original whitespace and remaining content
	buf.Write(data[insertPos:])

	// Write back to file
	if err := os.WriteFile(genesisPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write genesis file: %w", err)
	}

	return nil
}
