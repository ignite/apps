package cmd

import (
	"context"
	"os"

	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

func ExecuteHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

	return h.Run(
		ctx,
		hermes.WithArgs(cmd.Args...),
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}
