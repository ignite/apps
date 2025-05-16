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
func (r Querier) ValidateAppDetails(ctx context.Context, appFile, branch string) error {
	appBytes, err := os.ReadFile(appFile)
	if err != nil {
		return errors.Wrapf(err, "failed to get %s file content", appFile)
	}

	// load app entry from file.
	appEntry, err := AppFromFile(bytes.NewReader(appBytes))
	if err != nil {
		return err
	}

	// validate all JSON fields
	if err := appEntry.Validate(); err != nil {
		return err
	}

	// check if the name and ID is unique
	apps, err := r.List(ctx, branch)
	if err != nil {
		return err
	}

	if err := apps.CheckUnique(); err != nil {
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

	appYML, err := r.getAppsConfig(ctx, repo, branch)
	if err != nil {
		return err
	}

	// get the go.mod file from the app repository
	var goMod *modfile.File
	if repoOwner == "ignite" && repoName == "apps" {
		for id, info := range appYML.Apps {
			if !strings.EqualFold(id, appEntry.ID.String()) {
				continue
			}

			goMod, err = r.getGoMod(ctx, repo, path.Clean(info.Path), branch)
			if err != nil {
				return errors.Wrapf(err, "failed to get go.mod for app %s", id)
			}
			break
		}
		if goMod == nil {
			return errors.Errorf("oficial ignite app should be register into the %s file", appYMLFileName)
		}
	} else {
		goMod, err = r.getGoMod(ctx, repo, "", branch)
		if err != nil {
			return errors.Wrapf(err, "failed to get go.mod for app %s", appEntry.ID)
		}
	}

	cliVersion, err := findCLIVersion(goMod)
	if err != nil {
		return errors.Wrapf(err, "failed to find ignite version in go.mod for app %s", appEntry.ID)
	}

	if err := appEntry.Ignite.Verify(cliVersion); err != nil {
		return errors.Wrapf(err, "wrong ignite version %s for app %s: %s", cliVersion, appEntry.ID, appEntry.Ignite.String())
	}

	return nil
}
