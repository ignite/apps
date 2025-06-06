package cmd

import (
	"context"
	"os"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/wasm/pkg/config"
)

const (
	// flagSimulationGasLimit is the max gas to be used in a tx simulation call.
	// When not set the consensus max block gas is used instead.
	flagSimulationGasLimit = "simulation-gas-limit"
	// flagSmartQueryGasLimit is the max gas to be used in a smart query contract call.
	flagSmartQueryGasLimit = "query-gas-limit"
	// flagMemoryCacheSize in MiB not bytes.
	flagMemoryCacheSize = "memory-cache-size"
)

func ConfigHandler(ctx context.Context, cmd *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	flags := plugin.Flags(cmd.Flags)

	session := cliui.New(cliui.StartSpinnerWithText(statusAddingConfig))
	defer session.End()

	var (
		simulationGasLimit = getSimulationGasLimit(flags)
		smartQueryGasLimit = getSmartQueryGasLimit(flags)
		memoryCacheSize    = getMemoryCacheSize(flags)
	)

	c, err := newChain(ctx, api, chain.WithOutputer(session), chain.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	configTOML, err := c.ConfigTOMLPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(configTOML); os.IsNotExist(err) {
		return errors.Errorf("chain %s not initialized yet (%s)", c.Name(), c.AppPath())
	}

	if err := config.AddWasm(
		configTOML,
		config.WithSimulationGasLimit(simulationGasLimit),
		config.WithSmartQueryGasLimit(smartQueryGasLimit),
		config.WithMemoryCacheSize(memoryCacheSize),
	); err != nil {
		return err
	}
	session.Printf("\n🎉 CosmWasm config added at `%[1]v`.\n\n", configTOML)

	return nil
}
