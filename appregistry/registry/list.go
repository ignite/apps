package registry

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ignite/cli/v29/ignite/pkg/errors"

	"github.com/ignite/apps/appregistry/pkg/xgithub"
)

const (
	IgniteGitHubOrg  = "ignite"
	IgniteAppsRepo   = "apps"
	registryDir      = "_registry"
	igniteCLIPackage = "github.com/ignite/cli"
)

var appFormatRegex = regexp.MustCompile(`^([a-z]+\.[a-z]+\.[a-z]+\.json)$`)

type Querier struct {
	client *xgithub.Client
}

func NewRegistryQuerier(client *xgithub.Client) *Querier {
	return &Querier{client: client}
}

// List list apps from the ignite app appregistry/registry.
func (r *Querier) List(ctx context.Context, branch string) (Apps, error) {
	appsFiles, err := r.client.GetDirectoryFiles(
		ctx,
		IgniteGitHubOrg,
		IgniteAppsRepo,
		registryDir,
		xgithub.WithBranch(branch),
	)
	if err != nil {
		return nil, err
	}

	entries := make(Apps, 0)
	for _, file := range appsFiles {
		if !appFormatRegex.MatchString(strings.TrimPrefix(file, registryDir+"/")) {
			continue
		}

		entry, err := r.getRegistryEntry(ctx, file, branch)
		if err != nil {
			return nil, err
		}

		entries = append(entries, *entry)
	}

	return entries, nil
}

func (r *Querier) getRegistryEntry(ctx context.Context, fileName, branch string) (*App, error) {
	if branch == "" {
		branch = "main"
	}
	// here we do not use `GetFileContent` to avoid hitting the github api rate limit
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", IgniteGitHubOrg, IgniteAppsRepo, branch, fileName),
		nil,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build %s request", fileName)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s file content", fileName)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to get %s file content: %s", fileName, resp.Status)
	}

	namespace := namespaceFromFilePath(fileName)
	return AppFromFile(namespace, resp.Body)
}
