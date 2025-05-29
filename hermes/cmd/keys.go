package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

func KeysAddMnemonicHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	var (
		args          = cmd.Args
		flags         = plugin.Flags(cmd.Flags)
		hermesVersion = getVersion(flags)
		session       = cliui.New()
	)
	defer session.End()

	session.StartSpinner(fmt.Sprintf("Fetching hermes binary %s", hermesVersion))
	h, err := hermes.New(hermesVersion)
	if err != nil {
		return err
	}
	session.StopSpinner()

	return h.AddMnemonic(
		ctx,
		args[0],
		args[1],
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}

func KeysAddFileHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	var (
		args          = cmd.Args
		flags         = plugin.Flags(cmd.Flags)
		hermesVersion = getVersion(flags)
		session       = cliui.New()
	)
	defer session.End()

	session.StartSpinner(fmt.Sprintf("Fetching hermes binary %s", hermesVersion))
	h, err := hermes.New(hermesVersion)
	if err != nil {
		return err
	}
	session.StopSpinner()

	return h.AddKey(
		ctx,
		args[0],
		args[1],
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}

func KeysListHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	var (
		args          = cmd.Args
		flags         = plugin.Flags(cmd.Flags)
		hermesVersion = getVersion(flags)
		session       = cliui.New()
	)
	defer session.End()

	session.StartSpinner(fmt.Sprintf("Fetching hermes binary %s", hermesVersion))
	h, err := hermes.New(hermesVersion)
	if err != nil {
		return err
	}
	session.StopSpinner()

	return h.KeysList(
		ctx,
		args[0],
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}

func KeysDeleteHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	var (
		flags         = plugin.Flags(cmd.Flags)
		hermesVersion = getVersion(flags)
		session       = cliui.New()
	)
	defer session.End()

	session.StartSpinner(fmt.Sprintf("Fetching hermes binary %s", hermesVersion))
	h, err := hermes.New(hermesVersion)
	if err != nil {
		return err
	}
	session.StopSpinner()

	args := cmd.Args
	return h.DeleteKey(
		ctx,
		args[0],
		args[1],
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}
