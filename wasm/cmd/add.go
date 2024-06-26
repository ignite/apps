package cmd

import (
	"context"

	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/wasm/services/scaffolder"
)

func AddHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}

	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	var (
		simulationGasLimit = getSimulationGasLimit(flags)
		smartQueryGasLimit = getSmartQueryGasLimit(flags)
		memoryCacheSize    = getMemoryCacheSize(flags)
		wasmVersion        = getWasmVersion(flags)
	)

	wasmSemVer, err := semver.Parse(wasmVersion)
	if err != nil {
		return err
	}

	c, err := newChainWithHomeFlags(flags, chain.WithOutputer(session), chain.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(c, session)
	if err != nil {
		return err
	}

	sm, err := sc.AddWasm(
		ctx,
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
