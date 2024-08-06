package cmd

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/appregistry/pkg/tree"
	"github.com/ignite/apps/appregistry/pkg/xgithub"
	"github.com/ignite/apps/appregistry/registry"
)

const descriptionLimit = 75

// NewListCmd creates a new list command that lists all the ignite apps from the app registry.
func NewListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List all the ignite apps from the app registry",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			githubToken, _ := cmd.Flags().GetString(githubTokenFlag)

			session := cliui.New(cliui.StartSpinnerWithText("🔎 Searching for ignite apps on app registry..."))
			defer session.End()

			client := xgithub.NewClient(githubToken)
			registryQuerier := registry.NewRegistryQuerier(client)

			apps, err := registryQuerier.List(cmd.Context())
			if err != nil {
				return err
			}

			if len(apps) == 0 {
				session.Println("❌ No ignite application were found")
				return nil
			}

			session.StopSpinner()
			return session.Print(formatAppsTree(apps))
		},
	}

	return c
}

func formatAppsTree(entries []registry.AppEntry) string {
	b := &strings.Builder{}
	for _, entry := range entries {
		node := tree.NewNode(fmt.Sprintf(
			"%s : %s",
			entry.Name,
			limitTextLength(entry.Description, descriptionLimit),
		))
		node.AddChild(tree.NewNode(fmt.Sprintf(
			"📦 %s",
			entry.RepositoryURL,
		)))

		fmt.Fprint(b, node)
	}

	return b.String()
}

func limitTextLength(text string, limit int) string {
	if len(text) > limit {
		return text[:limit-3] + "..."
	}

	return text
}
