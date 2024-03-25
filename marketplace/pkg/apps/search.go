package apps

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/google/go-github/v56/github"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/marketplace/pkg/xgithub"
)

const (
	igniteAppTopic = "ignite-cli-app"
	appYMLFileName = "app.ignite.yml"
)

type (
	// AppRepository represents a GitHub repository with Ignite apps.
	AppRepository struct {
		PackageURL string
		Name       string
		Owner      string
		Stars      int
		UpdatedAt  time.Time
		Apps       []App
	}

	// App represents an Ignite app inside the repository.
	App struct {
		Name        string
		Description string
	}
)

// Search searches for repositories that have ignite app topic on GitHub given the query string and the minimum number of stars
// and then fetches the app.ignite.yml file from each repository and returns the list of repositories along with their apps.
func Search(ctx context.Context, client *xgithub.Client, query string, minStars uint) ([]AppRepository, error) {
	opts := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
	}
	repos, _, err := client.SearchRepositories(ctx, opts,
		xgithub.StringQuery(query),
		xgithub.LanguageQuery("go"),
		xgithub.TopicQuery(igniteAppTopic),
		xgithub.MinStarsQuery(int(minStars)))
	if err != nil {
		return nil, err
	}

	result := make([]AppRepository, 0, len(repos))
	for _, repo := range repos {
		apps, err := listApps(ctx, client, repo)

		var rateLimitErr *github.RateLimitError
		if errors.As(err, &rateLimitErr) {
			return nil, err
		} else if err != nil {
			// Ignore the repository since it doesn't have a valid app.ignite.yml file.
			continue
		}

		result = append(result, AppRepository{
			PackageURL: fmt.Sprintf("github.com/%s/%s", repo.GetOwner().GetLogin(), repo.GetName()),
			Name:       repo.GetName(),
			Owner:      repo.GetOwner().GetLogin(),
			Stars:      repo.GetStargazersCount(),
			UpdatedAt:  repo.GetPushedAt().Time,
			Apps:       apps,
		})
	}

	return result, nil
}

func listApps(ctx context.Context, client *xgithub.Client, repo *github.Repository) ([]App, error) {
	conf, err := getAppsConfig(ctx, client, repo)
	if err != nil {
		return nil, err
	}

	apps := make([]App, 0)
	for name, info := range conf.Apps {
		apps = append(apps, App{
			Name:        name,
			Description: info.Description,
		})
	}
	return apps, nil
}

func getAppsConfig(ctx context.Context, client *xgithub.Client, repo *github.Repository) (*plugin.AppsConfig, error) {
	data, err := client.GetFileContent(ctx, repo.GetOwner().GetLogin(), repo.GetName(), appYMLFileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s file content", appYMLFileName)
	}

	var conf plugin.AppsConfig
	if err := yaml.UnmarshalContext(ctx, data, &conf); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal %s file", appYMLFileName)
	}
	return &conf, nil
}
