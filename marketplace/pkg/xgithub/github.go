package xgithub

import (
	"context"

	"github.com/google/go-github/v56/github"
)

// SearchRepositories searches for repositories on GitHub given the query string and
// returns the list of repositories, the total number of results and an error.
func SearchRepositories(ctx context.Context, query, accToken string, page int) ([]*github.Repository, int, error) {
	client := github.NewClient(nil)
	if accToken != "" {
		client = client.WithAuthToken(accToken)
	}

	opts := &github.SearchOptions{Sort: "stars", Order: "desc"}
	if page > 0 {
		opts.Page = page
	}
	repos, _, err := client.Search.Repositories(ctx, query, &github.SearchOptions{Sort: "stars", Order: "desc"})
	if err != nil {
		return nil, 0, err
	}

	return repos.Repositories, *repos.Total, nil
}

// GetRepository gets the repository from GitHub given the repository name.
func GetRepository(ctx context.Context, owner, name, accToken string) (*github.Repository, error) {
	client := github.NewClient(nil)
	if accToken != "" {
		client = client.WithAuthToken(accToken)
	}

	repo, _, err := client.Repositories.Get(ctx, owner, name)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
