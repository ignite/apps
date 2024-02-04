package cmd

import (
	flag "github.com/spf13/pflag"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/wasmd/services/scaffolder"
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

	c.Flags().AddFlagSet(flagSetAppPath())
	return c
}

func flagSetAppPath() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringP(flagPath, "p", ".", "directory where the blockchain node is initialized")
	return fs
}

func scaffoldWasmHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusWasmInit))
	defer session.End()

	options := []scaffolder.WasmdOption{}
	sc, err := scaffolder.New(flagGetPath(cmd))
	if err != nil {
		return err
	}
	_, err = sc.InitWasmd(cmd.Context(), placeholder.New(), options...)
	return err
}
