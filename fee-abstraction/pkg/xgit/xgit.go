package xgit

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func HasVersion(versions semver.Versions, version semver.Version) bool {
	for _, v := range versions {
		if v.EQ(version) {
			return true
		}
	}
	return false
}

func FetchGitTags(repositoryURL string) (semver.Versions, error) {
	// Clone the repository
	path := filepath.Join(os.TempDir(), "fee-abstraction")
	if err := os.RemoveAll(path); err != nil {
		return nil, errors.Errorf("failed to clone repository: %s", err)
	}
	repo, err := git.PlainClone(
		path,
		false,
		&git.CloneOptions{URL: repositoryURL, Depth: 1},
	)
	if err != nil {
		return nil, errors.Errorf("failed to clone repository: %s", err)
	}

	// Get the repository's tags
	tags, err := repo.Tags()
	if err != nil {
		return nil, errors.Errorf("failed to get tags: %s", err)
	}

	// Iterate over tags and collect their names
	versions := make(semver.Versions, 0)
	err = tags.ForEach(func(tag *plumbing.Reference) error {
		versionName := strings.Replace(tag.Name().Short(), "v", "", 1)
		ver, err := semver.Parse(versionName)
		if err == nil {
			versions = append(versions, ver)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Errorf("failed to iterate over tags: %s", err)
	}
	semver.Sort(versions)
	return versions, nil
}
