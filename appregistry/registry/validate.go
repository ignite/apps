package registry

import (
	"bytes"
	"context"
	"os"
	"path"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

// ValidateAppDetails validates the details of an Ignite app repository.
func (r Querier) ValidateAppDetails(ctx context.Context, appFile string) error {
	appBytes, err := os.ReadFile(appFile)
	if err != nil {
		return errors.Wrapf(err, "failed to get %s file content", appFile)
	}

	appEntry, err := AppFromFile(bytes.NewReader(appBytes))
	if err != nil {
		return err
	}

	if err := appEntry.Validate(); err != nil {
		return err
	}

	repoOwner, repoName, err := validateRepoURL(appEntry.RepositoryURL.String())
	if err != nil {
		return err
	}

	repo, err := r.client.GetRepository(ctx, repoOwner, repoName)
	if err != nil {
		return err
	}

	appYML, err := r.getAppsConfig(ctx, repo)
	if err != nil {
		return err
	}

	var goMod *modfile.File
	if repoOwner == "ignite" && repoName == "apps" {
		for slug, info := range appYML.Apps {
			if !strings.EqualFold(slug, appEntry.Slug.String()) {
				continue
			}

			goMod, err = r.getGoMod(ctx, repo, path.Clean(info.Path))
			if err != nil {
				return errors.Wrapf(err, "failed to get go.mod for app %s", slug)
			}
			break
		}
		if goMod == nil {
			return errors.Errorf("oficial ignite app should be register into the %s file", appYMLFileName)
		}
	} else {
		goMod, err = r.getGoMod(ctx, repo, "")
		if err != nil {
			return errors.Wrapf(err, "failed to get go.mod for app %s", appEntry.Slug)
		}
	}

	cliVersion, err := findCLIVersion(goMod)
	if err != nil {
		return errors.Wrapf(err, "failed to find ignite version in go.mod for app %s", appEntry.Slug)
	}

	if err := appEntry.Ignite.Verify(cliVersion); err != nil {
		return errors.Wrapf(err, "failed to verify ignite version %s for app %s", cliVersion, appEntry.Slug)
	}

	return nil
}
