package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/official/marketplace/pkg/apps"
	"github.com/ignite/apps/official/marketplace/pkg/tree"
	"github.com/ignite/apps/official/marketplace/pkg/xgithub"
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
			githubToken, _ := cmd.Flags().GetString(githubTokenFlag)
			query, _ := cmd.Flags().GetString(queryFlag)
			minStars, _ := cmd.Flags().GetUint(minStarsFlag)

			session := cliui.New(cliui.StartSpinner())
			defer session.End()

			session.StartSpinner("🔎 Searching for ignite apps on GitHub...")
			client := xgithub.NewClient(githubToken)
			repos, err := apps.Search(cmd.Context(), client, query, minStars)
			if err != nil {
				return err
			}
			session.StopSpinner()

			if len(repos) < 1 {
				session.Println("❌ No ignite application were found")
				return nil
			}

			printRepoTree(session, repos)

			return nil
		},
	}

	c.Flags().StringP(queryFlag, "q", "", "Query string to search for")
	c.Flags().Uint(minStarsFlag, 10, "Minimum number of stars to search for")

	return c
}

func printRepoTree(sess *cliui.Session, repos []apps.AppRepository) {
	for i, repo := range repos {
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
				limitTextlength(app.Description, descriptionLimit),
			)))
		}
		sess.Print(node)

		if i < len(repos)-1 {
			sess.Println()
		}
	}
}

func humanizeInt(n int, unit string) string {
	value, suffix := humanize.ComputeSI(float64(n))
	return humanize.FtoaWithDigits(value, 1) + suffix + " " + unit
}

func limitTextlength(text string, limit int) string {
	if len(text) > limit {
		return text[:limit-3] + "..."
	}

	return text
}
