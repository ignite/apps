package cmd

import "github.com/spf13/cobra"

func NewRollkitAdd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add rollkit support",
		Long:  "Add rollkit support to your Cosmos SDK chain",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
}
