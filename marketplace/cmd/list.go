package cmd

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/marketplace/pkg/tree"
	"github.com/ignite/apps/marketplace/pkg/xgithub"
	"github.com/ignite/apps/marketplace/registry"
)

const descriptionLimit = 75

// NewList creates a new list command that searches all the ignite apps in GitHub.
func NewList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List all the ignite apps",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var githubToken, _ = cmd.Flags().GetString(githubTokenFlag)

			session := cliui.New(cliui.StartSpinnerWithText("🔎 Searching for ignite apps on GitHub..."))
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
	for i, entry := range entries {
		node := tree.NewNode(fmt.Sprintf(
			"%d. %-20s",
			i+1,
			entry.Name,
		))
		node.AddChild(tree.NewNode(fmt.Sprintf(
			"📖 %s",
			limitTextLength(entry.Description, descriptionLimit),
		)))
		node.AddChild(tree.NewNode(fmt.Sprintf(
			"📦 %s",
			entry.Repository.URL,
		)))

		fmt.Fprintln(b, node)
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func limitTextLength(text string, limit int) string {
	if len(text) > limit {
		return text[:limit-3] + "..."
	}

	return text
}
