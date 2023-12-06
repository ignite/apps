package cmd

import (
	"context"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v56/github"
	"github.com/ignite/apps/marketplace/pkg/xgithub"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/spf13/cobra"
)

const (
	minStarsFlag     = "min-stars"
	igniteAppTopic   = "ignite-cli-app"
	descriptionLimit = 50
)

// NewList creates a new list command that searches all the ignite apps in GitHub.
func NewList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list [query]",
		Short: "List all the ignite apps",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			minStars, _ := cmd.Flags().GetUint(minStarsFlag)

			session := cliui.New(cliui.StartSpinner())
			defer session.End()

			session.StartSpinner("🔎 Searching for ignite apps on GitHub...")
			repos, total, err := searchIgniteApps(cmd.Context(), query, minStars)
			if err != nil {
				return err
			}
			session.StopSpinner()

			session.Printf("🎉 Found %d results\n", total)

			if total > 0 {
				session.Println()
				printRepoList(session, repos)
			}

			return nil
		},
	}

	c.Flags().Uint(minStarsFlag, 10, "Minimum number of stars to search for")

	return c
}

func searchIgniteApps(ctx context.Context, query string, minStars uint) ([]*github.Repository, int, error) {
	client := xgithub.NewClient(githubToken)

	opts := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
	}
	repos, total, err := client.SearchRepositories(ctx, opts,
		xgithub.StringQuery(query),
		xgithub.LanguageQuery("go"),
		xgithub.TopicQuery(igniteAppTopic),
		xgithub.MinStarsQuery(int(minStars)))
	if err != nil {
		return nil, 0, err
	}

	return repos, total, nil
}

func printRepoList(sess *cliui.Session, repos []*github.Repository) {
	header := []string{"Name", "Description", "Stars ⭐️", "Updated At"}
	rows := make([][]string, 0, len(repos))
	for _, repo := range repos {
		rows = append(rows, []string{
			repo.GetFullName(),
			limitTextlength(repo.GetDescription(), descriptionLimit),
			humanize.SIWithDigits(float64(repo.GetStargazersCount()), 1, ""),
			humanize.Time(repo.GetPushedAt().Time),
		})
	}

	sess.PrintTable(header, rows...)
}

func limitTextlength(text string, limit int) string {
	if len(text) > limit {
		return text[:limit] + "..."
	}

	return text
}
