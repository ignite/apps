package cmd

import (
	"os"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/spf13/cobra"

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

// NewWasmConfig add wasm config to a chain.
func NewWasmConfig() *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "Add wasm config support",
		Args:  cobra.NoArgs,
		RunE:  wasmConfigExecuteHandler,
	}

	flagSetPath(c)
	flagSetHome(c)
	flagSetWasmConfigs(c)

	return c
}

func wasmConfigExecuteHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusAddingConfig))
	defer session.End()

	var (
		simulationGasLimit = getSimulationGasLimit(cmd)
		smartQueryGasLimit = getSmartQueryGasLimit(cmd)
		memoryCacheSize    = getMemoryCacheSize(cmd)
	)

	c, err := newChainWithHomeFlags(cmd, chain.WithOutputer(session), chain.CollectEvents(session.EventBus()))
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
