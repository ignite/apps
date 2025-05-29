package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

func ExecuteHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	var (
		flags   = plugin.Flags(cmd.Flags)
		session = cliui.New()
	)

	hermesVersion, err := getVersion(flags)
	if err != nil {
		return err
	}

	defer session.End()

	session.StartSpinner(fmt.Sprintf("Fetching hermes binary %s", hermesVersion))
	h, err := hermes.New(hermesVersion)
	if err != nil {
		return err
	}
	session.StopSpinner()

	return h.Run(
		ctx,
		hermes.WithArgs(cmd.Args...),
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}
