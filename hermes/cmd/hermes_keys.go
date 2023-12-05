package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

// NewHermesKeys manage the hermes relayer keys.
func NewHermesKeys() *cobra.Command {
	c := &cobra.Command{
		Use:   "keys",
		Short: "Manage the Hermes keys",
	}
	c.AddCommand(
		NewHermesKeysAddMnemonic(),
		NewHermesKeysAddFile(),
		NewHermesKeysList(),
		NewHermesKeysDelete(),
	)
	return c
}

// NewHermesKeysAddMnemonic add a hermes relayer mnemonic key.
func NewHermesKeysAddMnemonic() *cobra.Command {
	c := &cobra.Command{
		Use:   "add [chain-id] [mnemonic]",
		Short: "Add a new key from mnemonic to Hermes relayer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := hermes.New()
			if err != nil {
				return err
			}
			defer h.Cleanup()

			return h.AddMnemonic(
				cmd.Context(),
				args[0],
				args[1],
				hermes.WithStdOut(cmd.OutOrStdout()),
			)
		},
	}
	return c
}

// NewHermesKeysAddFile add a hermes relayer key file.
func NewHermesKeysAddFile() *cobra.Command {
	c := &cobra.Command{
		Use:   "file [chain-id] [filepath]",
		Short: "Add a new key from a key file to Hermes relayer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := hermes.New()
			if err != nil {
				return err
			}
			defer h.Cleanup()

			return h.AddKey(
				cmd.Context(),
				args[0],
				args[1],
				hermes.WithStdOut(cmd.OutOrStdout()),
			)
		},
	}
	return c
}

// NewHermesKeysList list hermes relayer keys.
func NewHermesKeysList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list [chain-id]",
		Short: "List Hermes relayer keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := hermes.New()
			if err != nil {
				return err
			}
			defer h.Cleanup()

			return h.KeysList(
				cmd.Context(),
				args[0],
				hermes.WithStdOut(cmd.OutOrStdout()),
			)
		},
	}
	return c
}

// NewHermesKeysDelete deletes a hermes relayer key.
func NewHermesKeysDelete() *cobra.Command {
	c := &cobra.Command{
		Use:   "delete [chain-id] [key-name]",
		Short: "Delete a key from Hermes relayer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := hermes.New()
			if err != nil {
				return err
			}
			defer h.Cleanup()

			return h.DeleteKey(
				cmd.Context(),
				args[0],
				args[1],
				hermes.WithStdOut(cmd.OutOrStdout()),
			)
		},
	}
	return c
}
