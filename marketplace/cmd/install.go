package cmd

import "github.com/spf13/cobra"

func NewInstallCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "install [name]",
		Short: "Install an ignite app by registry name",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {

			return nil
		},
	}

	return c
}
