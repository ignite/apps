package cmd

import (
	"github.com/ignite/apps/wasmd/services/scaffolder"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/spf13/cobra"
)

func NewWasmdInit() *cobra.Command {
	c := &cobra.Command{
		Use:     "init",
		Short:   "Import the wasm module to your app",
		Long:    "Add support for WebAssembly smart contracts to your blockchain",
		Args:    cobra.NoArgs,
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    scaffoldWasmHandler,
	}

	return c
}

func scaffoldWasmHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusWasmInit))
	defer session.End()

	var (
		appPath = "."
	)

	options := []scaffolder.WasmdOption{}
	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	if _, err = sc.InitWasmd(cmd.Context(), placeholder.New(), options...); err != nil {
		return err
	}
	return nil
}

//todo:: IBC checking need more investigation
