package cmd

import (
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/appregistry/pkg/xgithub"
	"github.com/ignite/apps/appregistry/registry"
)

// NewValidateCmd creates a new validate command that validates the Ignite application json.
func NewValidateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "validate [app file]",
		Aliases: []string{"v"},
		Short:   "Validate the ignite application json",
		Args:    cobra.ExactArgs(1),
		RunE:    validateHandler,
	}

	c.Flags().StringP(flagBranch, "b", "", "The app branch to use (default: main)")

	return c
}

func validateHandler(cmd *cobra.Command, args []string) error {
	var (
		githubToken, _ = cmd.Flags().GetString(flagGithubToken)
		branch, _      = cmd.Flags().GetString(flagBranch)
	)

	session := cliui.New(cliui.StartSpinnerWithText("🔎 Fetching repository details from GitHub..."))
	defer session.End()

	client := xgithub.NewClient(githubToken)
	registryQuerier := registry.NewRegistryQuerier(client)

	absPath, err := filepath.Abs(args[0])
	if err != nil {
		return errors.Wrapf(err, "failed to get absolute path for %s", args[0])
	}
	if err := registryQuerier.ValidateAppDetails(cmd.Context(), absPath, branch); err != nil {
		return err
	}

	session.StopSpinner()
	return session.Printf("🚀 valid %s file\n", args[0])
}
