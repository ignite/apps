package cmd

import "github.com/spf13/cobra"

func NewWasmdSC() *cobra.Command {
	c := &cobra.Command{
		Use:   "sc",
		Short: "Import the wasm module to your app",
		Long:  "Add support for WebAssembly smart contracts to your blockchain",
		Args:  cobra.NoArgs,
		RunE:  scaffoldWasmSCHandler,
	}

	// flagSetPath(c)
	//	TODO :: add subcommands like new, test, deploy, ...
	return c
}

func scaffoldWasmSCHandler(cmd *cobra.Command, _ []string) error {
	return nil
}
