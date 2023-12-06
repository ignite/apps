package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v56/github"
	"github.com/gookit/color"
	"github.com/ignite/apps/marketplace/pkg/xgithub"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/spf13/cobra"
)

// NewInfo creates a new info command that shows the details of an ignite app.
func NewInfo() *cobra.Command {
	return &cobra.Command{
		Use:   "info [app-name]",
		Short: "Show the details of an ignite app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			repoOwner, repoName, err := validateAppName(appName)
			if err != nil {
				return err
			}

			session := cliui.New(cliui.StartSpinner())
			defer session.End()

			session.StartSpinner("Fetching app details from GitHub...")
			repo, err := getRepo(cmd.Context(), repoOwner, repoName)
			if err != nil {
				return err
			}
			session.StopSpinner()

			if !isIgniteApp(repo) {
				return fmt.Errorf("the repository is not an ignite app")
			}

			printRepoDetails(session, repo)

			return nil
		},
	}
}

func validateAppName(appName string) (owner, name string, err error) {
	appName = strings.TrimPrefix(appName, "github.com/")
	appNameParts := strings.Split(appName, "/")
	if len(appNameParts) != 2 {
		return "", "", fmt.Errorf("invalid app name: %s", appName)
	}
	owner = appNameParts[0]
	name = appNameParts[1]

	return owner, name, nil
}

func getRepo(ctx context.Context, owner, name string) (*github.Repository, error) {
	client := xgithub.NewClient(githubToken)
	return client.GetRepository(ctx, owner, name)
}

func isIgniteApp(repo *github.Repository) bool {
	for _, topic := range repo.Topics {
		if topic == "ignite-app" {
			return true
		}
	}

	return false
}

func printRepoDetails(sess *cliui.Session, repo *github.Repository) {
	sess.Println("Name: ", repo.GetName())
	sess.Println("Owner: ", repo.GetOwner().GetLogin())
	sess.Println("Description: ", repo.GetDescription())
	sess.Println("Stars ⭐️: ", repo.GetStargazersCount())
	sess.Println("License: ", repo.GetLicense().GetName())
	sess.Printf("Updated At: %s (%s)\n", repo.GetUpdatedAt().Format(time.DateTime), humanize.Time(repo.GetUpdatedAt().Time))
	sess.Println("URL: ", repo.GetHTMLURL())
	sess.Println()
	sess.Printf("🚀 Install via: %s\n", color.Green.Sprintf("ignite app install github.com/%s", repo.GetFullName()))
}
