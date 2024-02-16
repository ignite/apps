package cmd

import (
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
	c.Flags().Uint64(flagSimulationGasLimit, 0, "the max gas to be used in a tx simulation call. When not set the consensus max block gas is used instead")
	c.Flags().Uint64(flagSmartQueryGasLimit, 3_000_000, "the max gas to be used in a smart query contract call")
	c.Flags().Uint32(flagMemoryCacheSize, 100, "memory cache size in MiB not bytes")

	return c
}

func wasmAddExecuteHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	var (
		simulationGasLimit, _ = cmd.Flags().GetUint64(flagSimulationGasLimit)
		smartQueryGasLimit, _ = cmd.Flags().GetUint64(flagSmartQueryGasLimit)
		memoryCacheSize, _    = cmd.Flags().GetUint32(flagMemoryCacheSize)
	)

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
