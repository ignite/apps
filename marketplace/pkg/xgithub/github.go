package xgithub

import (
	"context"
	"strconv"

	"github.com/google/go-github/v56/github"
	"github.com/pkg/errors"
)

// Client is a wrapper around the GitHub client so that it can be used as
// an implementation of interface later on.
type Client struct {
	gc *github.Client
}

// NewClient returns a new GitHub client.
func NewClient(accessToken string) *Client {
	gc := github.NewClient(nil)
	if accessToken != "" {
		gc = gc.WithAuthToken(accessToken)
	}

	return &Client{gc: gc}
}

type Query struct {
	qualifier string
	op        string
	value     string
}

func (q Query) String() string {
	return q.qualifier + q.op + q.value
}

// StringQuery returns a Query that matches the given string.
func StringQuery(str string) Query {
	return Query{qualifier: str}
}

// TopicQuery returns a Query that matches the given topic.
func TopicQuery(topic string) Query {
	return Query{qualifier: "topic", op: ":", value: topic}
}

// LanguageQuery returns a Query that matches the given language.
func LanguageQuery(lang string) Query {
	return Query{qualifier: "language", op: ":", value: lang}
}

// MinStarsQuery returns a Query that matches the given minimum number of stars.
func MinStarsQuery(stars int) Query {
	return Query{qualifier: "stars", op: ":>=", value: strconv.Itoa(stars)}
}

// SearchRepositories searches for repositories on GitHub given the query string and
// returns the list of repositories, the total number of results and an error.
func (c *Client) SearchRepositories(ctx context.Context, opts *github.SearchOptions, queries ...Query) ([]*github.Repository, int, error) {
	q := joinQueries(queries...)
	repos, _, err := c.gc.Search.Repositories(ctx, q, opts)
	if err != nil {
		return nil, 0, err
	}

	return repos.Repositories, *repos.Total, nil
}

func joinQueries(queries ...Query) string {
	var q string
	for i, query := range queries {
		q += query.String()

		if i < len(queries)-1 {
			q += " "
		}
	}

	return q
}

// GetRepository gets the repository from GitHub given the repository name.
func (c *Client) GetRepository(ctx context.Context, owner, name string) (*github.Repository, error) {
	repo, _, err := c.gc.Repositories.Get(ctx, owner, name)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// GetFileContent gets the content of the file from GitHub given the repository name and the file path.
func (c *Client) GetFileContent(ctx context.Context, owner, repo, path string) ([]byte, error) {
	file, _, _, err := c.gc.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, err
	}

	s, err := file.GetContent()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file content")
	}

	return []byte(s), nil
}
