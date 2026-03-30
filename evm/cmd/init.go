package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

const eighteenDecimals = 18

func InitHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
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

	c, err := chain.New(absPath, chain.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	igniteConfig, err := c.Config()
	if err != nil {
		return err
	}

	denom := igniteConfig.DefaultDenom
	if denom == "" {
		denom = sdk.DefaultBondDenom
	}

	if err := c.Init(ctx, chain.InitArgsAll); err != nil {
		return err
	}

	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}

	genesis, err := genutiltypes.AppGenesisFromFile(genesisPath)
	if err != nil {
		return err
	}

	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesis.AppState, &appState); err != nil {
		return err
	}

	bankGenesisBz, ok := appState[banktypes.ModuleName]
	if !ok {
		return fmt.Errorf("bank module not found in genesis app state")
	}

	var bankGenesis banktypes.GenesisState
	if err := json.Unmarshal(bankGenesisBz, &bankGenesis); err != nil {
		return err
	}

	displayDenom := denom
	if strings.HasPrefix(denom, "u") && len(denom) > 1 {
		displayDenom = denom[1:]
	}

	bankGenesis.DenomMetadata = []banktypes.Metadata{
		{
			Description: fmt.Sprintf("Native 18-decimal denom metadata for %s EVM chain", denom),
			Base:        denom,
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: denom, Exponent: 0},
				{Denom: displayDenom, Exponent: eighteenDecimals},
			},
			Name:    displayDenom,
			Symbol:  strings.ToUpper(displayDenom),
			Display: displayDenom,
		},
	}

	bankGenesisBz, err = json.Marshal(bankGenesis)
	if err != nil {
		return err
	}
	appState[banktypes.ModuleName] = bankGenesisBz

	genesis.AppState, err = json.Marshal(appState)
	if err != nil {
		return err
	}

	if err := genesis.SaveAs(genesisPath); err != nil {
		return err
	}

	home, err := c.Home()
	if err != nil {
		return err
	}

	return session.Printf("🗃 Initialized. Checkout your chain's home (data) directory: %s\n", colors.Info(home))
}
