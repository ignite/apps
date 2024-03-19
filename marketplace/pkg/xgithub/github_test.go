package xgithub

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v56/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SearchRepositories(t *testing.T) {
	var (
		require      = require.New(t)
		assert       = assert.New(t)
		searchResult = &github.RepositoriesSearchResult{
			Total: github.Int(2),
			Repositories: []*github.Repository{
				{
					FullName:        github.String("igniteapps/ignite-cli-app-2"),
					Description:     github.String("This is a test ignite app 2"),
					StargazersCount: github.Int(20),
					Topics:          []string{"ignite-cli-app"},
				},
				{
					FullName:        github.String("igniteapps/ignite-cli-app-1"),
					Description:     github.String("This is a test ignite app 1"),
					StargazersCount: github.Int(10),
					Topics:          []string{"ignite-cli-app"},
				},
			},
		}
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v3/search/repositories", r.URL.Path)
		assert.Equal("test topic:ignite-cli-app stars:>=10", r.FormValue("q"))
		assert.Equal("stars", r.FormValue("sort"))

		err := json.NewEncoder(w).Encode(searchResult)
		require.NoError(err)
	}))
	defer ts.Close()

	gc := github.NewClient(nil)
	gc, err := gc.WithEnterpriseURLs(ts.URL, ts.URL)
	require.NoError(err)

	client := &Client{gc: gc}
	repos, total, err := client.SearchRepositories(
		context.Background(),
		&github.SearchOptions{Sort: "stars", Order: "desc"},
		StringQuery("test"), TopicQuery("ignite-cli-app"), MinStarsQuery(10),
	)
	require.NoError(err)
	assert.Equal(*searchResult.Total, total)
	assert.Equal(searchResult.Repositories, repos)
}

func TestClient_GetRepository(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
		repo    = &github.Repository{
			Name: github.String("ignite-cli-app"),
			Owner: &github.User{
				Login: github.String("igniteapps"),
			},
			Description:     github.String("This is a test ignite app"),
			StargazersCount: github.Int(10),
			License: &github.License{
				Name: github.String("MIT"),
			},
			HTMLURL: github.String("https://github.com/igniteapps/ignite-cli-app"),
		}
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v3/repos/igniteapps/ignite-cli-app", r.URL.Path)

		err := json.NewEncoder(w).Encode(repo)
		require.NoError(err)
	}))
	defer ts.Close()

	gc := github.NewClient(nil)
	gc, err := gc.WithEnterpriseURLs(ts.URL, ts.URL)
	require.NoError(err)

	client := &Client{gc: gc}
	repo, err = client.GetRepository(context.Background(), "igniteapps", "ignite-cli-app")
	require.NoError(err)
	assert.Equal(repo, repo)
}
