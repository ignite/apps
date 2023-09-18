package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

// NewHermesExecute execute hermes relayer commands.
func NewHermesExecute() *cobra.Command {
	c := &cobra.Command{
		Use:   "exec [args...]",
		Short: "Execute a hermes raw command",
		Args:  cobra.MinimumNArgs(1),
		RunE:  hermesExecuteHandler,
	}

	return c
}

func hermesExecuteHandler(cmd *cobra.Command, args []string) error {
	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()
	return h.Run(cmd.Context(), os.Stdout, os.Stderr, "", args...)
}
