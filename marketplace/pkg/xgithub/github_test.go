package xgithub_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v56/github"
	"github.com/stretchr/testify/require"

	"github.com/ignite/apps/marketplace/pkg/xgithub"
)

func TestClient_GetDirectoryFiles(t *testing.T) {
	gc := xgithub.NewClient("")

	files, err := gc.GetDirectoryFiles(context.Background(), "ignite", "apps", "_registry")
	require.NoError(t, err)

	require.Contains(t, files, "_registry/README.md")
	require.Contains(t, files, "_registry/registry.json")
	require.GreaterOrEqual(t, len(files), 5)
}

func TestClient_GetRepository(t *testing.T) {
	var (
		require = require.New(t)
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
		require.Equal("/api/v3/repos/igniteapps/ignite-cli-app", r.URL.Path)

		err := json.NewEncoder(w).Encode(repo)
		require.NoError(err)
	}))
	defer ts.Close()

	gc := github.NewClient(nil)
	gc, err := gc.WithEnterpriseURLs(ts.URL, ts.URL)
	require.NoError(err)

	client := &xgithub.Client{GithubClient: gc}
	repo, err = client.GetRepository(context.Background(), "igniteapps", "ignite-cli-app")
	require.NoError(err)
	require.Equal(repo, repo)
}
