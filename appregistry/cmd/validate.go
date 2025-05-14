package cmd

import (
	"fmt"
	"strconv"
	"text/tabwriter"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/appregistry/pkg/xgithub"
	"github.com/ignite/apps/appregistry/registry"
)

// NewValidateCmd creates a new validate command that validates the ignite application yaml.
func NewValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "validate [app name]",
		Aliases: []string{"v"},
		Short:   "Validate the ignite application yaml",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			githubToken, _ := cmd.Flags().GetString(githubTokenFlag)

			session := cliui.New(cliui.StartSpinnerWithText("🔎 Fetching repository details from GitHub..."))
			defer session.End()

			client := xgithub.NewClient(githubToken)
			registryQuerier := registry.NewRegistryQuerier(client)

			appDetails, err := registryQuerier.ValidateApp(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			session.StopSpinner()

			w := &tabwriter.Writer{}
			printItem := func(s string, v interface{}) {
				fmt.Fprintf(w, "\t%s:\t%v\n", s, v)
			}
			w.Init(cmd.OutOrStdout(), 0, 8, 0, '\t', 0)

			printItem("Name", appDetails.App.Name)
			printItem("Description", appDetails.App.Description)
			printItem("Stars", strconv.Itoa(appDetails.Stars))
			printItem("Go version", appDetails.App.GoVersion)
			printItem("Ignite version", appDetails.App.IgniteVersion)
			printItem("Documentation", appDetails.App.DocumentationURL)
			printItem("Repository", linkStyle.Render(appDetails.URL))

			fmt.Fprintln(w, installationStyle.Render(fmt.Sprintf(
				"🚀 Install via: %s", commandStyle.Render(fmt.Sprintf("ignite app -g install %s", appDetails.App.PackageURL)),
			)))

			w.Flush()

			return nil
		},
	}
}
