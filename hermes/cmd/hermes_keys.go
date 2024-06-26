package cmd

import (
	"context"
	"os"

	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

func KeysAddMnemonicHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	args := cmd.Args
	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

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
	args := cmd.Args
	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

	return h.AddKey(
		ctx,
		args[0],
		args[1],
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}

func KeysList(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	args := cmd.Args
	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

	return h.KeysList(
		ctx,
		args[0],
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}

func NewHermesKeysDelete(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

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
