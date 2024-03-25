package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/marketplace/pkg/apps"
	"github.com/ignite/apps/marketplace/pkg/tree"
	"github.com/ignite/apps/marketplace/pkg/xgithub"
)

const (
	minStarsFlag     = "min-stars"
	queryFlag        = "query"
	descriptionLimit = 75
)

var (
	starsCountStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	updatedAtStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

// NewList creates a new list command that searches all the ignite apps in GitHub.
func NewList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List all the ignite apps",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				githubToken, _ = cmd.Flags().GetString(githubTokenFlag)
				query, _       = cmd.Flags().GetString(queryFlag)
				minStars, _    = cmd.Flags().GetUint(minStarsFlag)
			)

			session := cliui.New(cliui.StartSpinnerWithText("🔎 Searching for ignite apps on GitHub..."))
			defer session.End()

			client := xgithub.NewClient(githubToken)
			repos, err := apps.Search(cmd.Context(), client, query, minStars)
			if err != nil {
				return err
			}

			if len(repos) < 1 {
				return fmt.Errorf("❌ No ignite application were found")
			}

			session.StopSpinner()
			return session.Print(formatRepoTree(repos))
		},
	}

	c.Flags().StringP(queryFlag, "q", "", "Query string to search for")
	c.Flags().Uint(minStarsFlag, 0, "Minimum number of stars to search for")

	return c
}

func formatRepoTree(repos []apps.AppRepository) string {
	b := &strings.Builder{}
	for _, repo := range repos {
		node := tree.NewNode(fmt.Sprintf(
			"📦 %-50s %s %s",
			repo.PackageURL,
			starsCountStyle.Render(humanizeInt(repo.Stars, "⭐️")),
			updatedAtStyle.Render("("+humanize.Time(repo.UpdatedAt)+")"),
		))
		node.AddChild(nil) // Add a nil child to add a line break.
		for _, app := range repo.Apps {
			node.AddChild(tree.NewNode(fmt.Sprintf(
				"🔥 %-20s %s",
				app.Name,
				limitTextLength(app.Description, descriptionLimit),
			)))
		}
		fmt.Fprintln(b, node)
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func humanizeInt(n int, unit string) string {
	value, suffix := humanize.ComputeSI(float64(n))
	return humanize.FtoaWithDigits(value, 1) + suffix + " " + unit
}

func limitTextLength(text string, limit int) string {
	if len(text) > limit {
		return text[:limit-3] + "..."
	}

	return text
}
