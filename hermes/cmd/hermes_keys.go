package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"relayer/pkg/hermes"
)

// NewHermesKeys manage the hermes relayer keys.
func NewHermesKeys() *cobra.Command {
	c := &cobra.Command{
		Use:   "keys",
		Short: "manage the Hermes keys",
	}
	c.AddCommand(
		NewHermesKeyAddMnemonic(),
		NewHermesKeyAddFile(),
	)
	return c
}

// NewHermesKeyAddMnemonic add a hermes relayer mnemonic key.
func NewHermesKeyAddMnemonic() *cobra.Command {
	c := &cobra.Command{
		Use:   "add [chain-id] [mnemonic]",
		Short: "add a new key from mnemonic to Hermes relayer",
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
				hermes.WithStdOut(os.Stdout),
			)
		},
	}
	return c
}

// NewHermesKeyAddFile add a hermes relayer key file.
func NewHermesKeyAddFile() *cobra.Command {
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
				hermes.WithStdOut(os.Stdout),
			)
		},
	}
	return c
}
