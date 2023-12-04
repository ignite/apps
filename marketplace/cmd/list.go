package cmd

import (
	"context"
	"fmt"

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
			minStars, _ := cmd.Flags().GetUint(minStarsFlag)
			query := fmt.Sprintf("topic:%s language:go stars:>=%d", igniteAppTopic, minStars)
			if len(args) > 0 {
				query = args[0] + " " + query
			}

			sess := cliui.New()
			defer sess.End()

			sess.StartSpinner("🔎 Searching for ignite apps on GitHub...")
			repos, total, err := searchIgniteApps(cmd.Context(), query, githubAccessToken)
			if err != nil {
				return err
			}
			sess.StopSpinner()

			sess.Printf("🎉 Found %d results\n", total)

			if total > 0 {
				sess.Println()
				printRepoList(sess, repos)
			}

			return nil
		},
	}

	c.LocalFlags().Uint(minStarsFlag, 10, "Minimum number of stars to search for")

	return c
}

func searchIgniteApps(ctx context.Context, query, accToken string) ([]*github.Repository, int, error) {
	repos, total, err := xgithub.SearchRepositories(ctx, query, accToken, 0)
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
