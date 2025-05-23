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
	return &cobra.Command{
		Use:   "list",
		Short: "List all the ignite apps from the app registry",
		Args:  cobra.NoArgs,
		RunE:  listHandler,
	}
}

func listHandler(cmd *cobra.Command, _ []string) error {
	var (
		githubToken = getGitHubToken(cmd)
		branch      = getBranchFlag(cmd)
	)

	session := cliui.New(cliui.StartSpinnerWithText("🔎 Searching for ignite apps on app registry..."))
	defer session.End()

	client := xgithub.NewClient(githubToken)
	registryQuerier := registry.NewRegistryQuerier(client)

	apps, err := registryQuerier.List(cmd.Context(), branch)
	if err != nil {
		return err
	}

	if len(apps) == 0 {
		session.Println("❌ No ignite application were found")
		return nil
	}

	session.StopSpinner()
	return session.Print(formatAppsTree(apps))
}

func formatAppsTree(entries []registry.App) string {
	b := &strings.Builder{}
	for _, entry := range entries {
		node := tree.NewNode(fmt.Sprintf(
			"%s (id: %s): %s",
			entry.Name,
			entry.AppID,
			limitTextLength(entry.Description.String(), descriptionLimit),
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
