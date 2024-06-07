package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/errors"

	"github.com/ignite/apps/marketplace/pkg/xgithub"
)

const (
	igniteGitHubOrg  = "ignite"
	igniteAppsRepo   = "apps"
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

// List list apps from the ignite app marketplace/registry.
func (r *Querier) List(ctx context.Context) ([]AppEntry, error) {
	appsFiles, err := r.client.GetDirectoryFiles(ctx, igniteGitHubOrg, igniteAppsRepo, registryDir)
	if err != nil {
		return nil, err
	}

	entries := make([]AppEntry, 0)
	for _, file := range appsFiles {
		if !appFormatRegex.MatchString(strings.TrimPrefix(file, registryDir+"/")) {
			continue
		}

		entry, err := r.getRegistryEntry(file)
		if err != nil {
			return nil, err
		}

		entries = append(entries, *entry)
	}

	return entries, nil
}

func (r *Querier) getRegistryEntry(fileName string) (*AppEntry, error) {
	// here we do not use `GetFileContent` to avoid hitting the github api rate limit
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/julien/marketplace/%s", igniteGitHubOrg, igniteAppsRepo, fileName))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s file content", fileName)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s file content", fileName)
	}

	var entry *AppEntry
	if err := json.Unmarshal(body, &entry); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal %s file", fileName)
	}

	return entry, nil
}
