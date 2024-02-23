package cmd

import (
	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/wasm/services/scaffolder"
)

// NewWasmAdd add wasm integration to a chain.
func NewWasmAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add",
		Short: "Add wasm support",
		Args:  cobra.NoArgs,
		RunE:  wasmAddExecuteHandler,
	}

	flagSetPath(c)
	flagSetHome(c)
	flagSetWasmConfigs(c)
	flagSetWasmVersion(c)

	return c
}

func wasmAddExecuteHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	var (
		simulationGasLimit = getSimulationGasLimit(cmd)
		smartQueryGasLimit = getSmartQueryGasLimit(cmd)
		memoryCacheSize    = getMemoryCacheSize(cmd)
		wasmVersion        = getWasmVersion(cmd)
	)

	wasmSemVer, err := semver.Parse(wasmVersion)
	if err != nil {
		return err
	}

	c, err := newChainWithHomeFlags(cmd, chain.WithOutputer(session), chain.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(c, session)
	if err != nil {
		return err
	}

	sm, err := sc.AddWasm(
		cmd.Context(),
		placeholder.New(),
		scaffolder.WithWasmVersion(wasmSemVer),
		scaffolder.WithSimulationGasLimit(simulationGasLimit),
		scaffolder.WithSmartQueryGasLimit(smartQueryGasLimit),
		scaffolder.WithMemoryCacheSize(memoryCacheSize),
	)
	if err != nil {
		return err
	}

	modificationsStr, err := sourceModificationToString(sm)
	if err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 CosmWasm added (`%[1]v`).\n\n", c.AppPath())

	return nil
}
