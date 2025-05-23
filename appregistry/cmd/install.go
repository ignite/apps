package cmd

import (
	"os"

	"github.com/spf13/cobra"

	ignitecmd "github.com/ignite/cli/v28/ignite/cmd"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"

	"github.com/ignite/apps/appregistry/pkg/xgithub"
	"github.com/ignite/apps/appregistry/registry"
)

func NewInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install [app id]",
		Short: "Install an ignite app by app id",
		Args:  cobra.ExactArgs(1),
		RunE:  installHandler,
	}
}

func installHandler(cmd *cobra.Command, args []string) error {
	var (
		githubToken = getGitHubToken(cmd)
		branch      = getBranchFlag(cmd)
	)

	session := cliui.New(cliui.WithStdout(os.Stdout))
	defer session.End()

	client := xgithub.NewClient(githubToken)
	registryQuerier := registry.NewRegistryQuerier(client)

	appDetails, err := registryQuerier.GetAppDetails(cmd.Context(), args[0], branch)
	if err != nil {
		return err
	}

	// here we are using the ignite app install command to install the app
	// we do this in order to not duplicate logic.
	igniteAppInstallCmd := ignitecmd.NewAppInstall()
	igniteAppInstallCmd.SetArgs([]string{"-g", appDetails.App.PackageURL})

	return igniteAppInstallCmd.ExecuteContext(cmd.Context())
}
