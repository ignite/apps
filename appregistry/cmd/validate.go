package cmd

import (
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/appregistry/pkg/xgithub"
	"github.com/ignite/apps/appregistry/registry"
)

// NewValidateCmd creates a new validate command that validates the Ignite application json.
func NewValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "validate [app file]",
		Aliases: []string{"v"},
		Short:   "Validate the ignite application json",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			githubToken, _ := cmd.Flags().GetString(githubTokenFlag)

			session := cliui.New(cliui.StartSpinnerWithText("🔎 Fetching repository details from GitHub..."))
			defer session.End()

			client := xgithub.NewClient(githubToken)
			registryQuerier := registry.NewRegistryQuerier(client)

			if err := registryQuerier.ValidateAppDetails(cmd.Context(), args[0]); err != nil {
				return err
			}

			session.StopSpinner()

			return session.Printf("🚀 valid %s file", args[0])
		},
	}
}
