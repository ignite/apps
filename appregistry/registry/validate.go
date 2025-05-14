package registry

import (
	"bytes"
	"context"
	"os"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

// ValidateAppDetails validates the details of an Ignite app repository.
func (r Querier) ValidateAppDetails(ctx context.Context, AppFile string) error {
	appBytes, err := os.ReadFile(AppFile)
	if err != nil {
		return errors.Wrapf(err, "failed to get %s file content", AppFile)
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

	goMod, err := r.getGoMod(ctx, repo, "")
	if err != nil {
		return errors.Wrapf(err, "failed to get go.mod for app %s", appEntry.Name)
	}

	cliVersion, err := findCLIVersion(goMod)
	if err != nil {
		return errors.Wrapf(err, "failed to find ignite version in go.mod for app %s", appEntry.Name)
	}

	if err := appEntry.Ignite.Verify(cliVersion); err != nil {
		return errors.Wrapf(err, "failed to verify ignite version %s for app %s", cliVersion, appEntry.Name)
	}

	return nil
}
