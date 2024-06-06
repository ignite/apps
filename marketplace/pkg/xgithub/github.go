package xgithub

import (
	"context"

	"github.com/google/go-github/v56/github"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

// Client is a wrapper around GitHub client.
type Client struct {
	GithubClient *github.Client
}

// NewClient returns a new GitHub client.
func NewClient(accessToken string) *Client {
	gc := github.NewClient(nil)
	if accessToken != "" {
		gc = gc.WithAuthToken(accessToken)
	}

	return &Client{GithubClient: gc}
}

// GetRepository gets the repository from GitHub given the repository name.
func (c *Client) GetRepository(ctx context.Context, owner, name string) (*github.Repository, error) {
	repo, _, err := c.GithubClient.Repositories.Get(ctx, owner, name)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// GetDirectoryFiles lists the files paths in the directory from GitHub given the repository name and the directory path.
func (c *Client) GetDirectoryFiles(ctx context.Context, owner, repo, path string) ([]string, error) {
	_, dir, _, err := c.GithubClient.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, err
	}

	var filesPaths []string
	for _, f := range dir {
		filesPaths = append(filesPaths, f.GetPath())
	}

	return filesPaths, nil
}

// GetFileContent gets the content of the file from GitHub given the repository name and the file path.
func (c *Client) GetFileContent(ctx context.Context, owner, repo, path string) ([]byte, error) {
	file, _, _, err := c.GithubClient.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, err
	}

	s, err := file.GetContent()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file content")
	}

	return []byte(s), nil
}
