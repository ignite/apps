package cmd

import (
	"fmt"
	"strconv"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/appregistry/pkg/xgithub"
	"github.com/ignite/apps/appregistry/registry"
)

var (
	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Underline(true)

	installationStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("9"))

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)
)

// NewDetailsCmd creates a new details command that shows the details of an ignite application repository.
func NewDetailsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "details [app id]",
		Aliases: []string{"info"},
		Short:   "Show the details of an ignite application repository",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				githubToken, _ = cmd.Flags().GetString(flagGithubToken)
				branch, _      = cmd.Flags().GetString(flagBranch)
			)

			session := cliui.New(cliui.StartSpinnerWithText("🔎 Fetching repository details from GitHub..."))
			defer session.End()

			client := xgithub.NewClient(githubToken)
			registryQuerier := registry.NewRegistryQuerier(client)

			appDetails, err := registryQuerier.GetAppDetails(cmd.Context(), args[0], branch)
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
			printItem("ID", appDetails.App.ID)
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
