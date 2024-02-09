package cmd

import (
	"wasm/services/scaffolder"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/spf13/cobra"
)

// NewWasm creates a new wasm command that holds
// some other sub commands related to CosmWasm.
func NewWasm() *cobra.Command {
	c := &cobra.Command{
		Use:           "wasm [command]",
		Short:         "Ignite wasm integration",
		Aliases:       []string{"w"},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// add sub commands.
	c.AddCommand(
		NewWasmAdd(),
	)
	return c
}

// NewWasmAdd add wasm integration to a chain.
func NewWasmAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add",
		Short: "Add wasm support",
		Args:  cobra.NoArgs,
		RunE:  wasmAddExecuteHandler,
	}

	flagSetPath(c)

	return c
}

func wasmAddExecuteHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	appPath := flagGetPath(cmd)
	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}

	sm, err := sc.AddWasm(cmd.Context(), placeholder.New())
	if err != nil {
		return err
	}

	modificationsStr, err := sourceModificationToString(sm)
	if err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 CosmWasm added (`%[1]v`).\n\n", appPath)

	return nil
}
