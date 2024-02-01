package apps

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/go-github/v56/github"
	"github.com/ignite/apps/official/marketplace/pkg/xgithub"
	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
)

const igniteCLIPackage = "github.com/ignite/cli"

// AppRepositoryDetails represents the details of an Ignite app repository.
type AppRepositoryDetails struct {
	PackageURL  string
	Name        string
	Owner       string
	Description string
	Tags        []string
	Stars       int
	License     string
	UpdatedAt   time.Time
	URL         string
	Apps        []AppDetails
}

// AppDetails represents the details of an Ignite app.
type AppDetails struct {
	Name          string
	Description   string
	Path          string
	GoVersion     string
	IgniteVersion string
}

// GetRepositoryDetails returns the details of an Ignite app repository.
func GetRepositoryDetails(ctx context.Context, client *xgithub.Client, pkgURL string) (*AppRepositoryDetails, error) {
	repoOwner, repoName, err := validatePackageURL(pkgURL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid package URL")
	}

	repo, err := client.GetRepository(ctx, repoOwner, repoName)
	if err != nil {
		return nil, err
	}

	appYML, err := getAppsConfig(ctx, client, repo)
	if err != nil {
		return nil, err
	}

	result := &AppRepositoryDetails{
		PackageURL:  pkgURL,
		Name:        repo.GetName(),
		Owner:       repo.GetOwner().GetLogin(),
		Description: repo.GetDescription(),
		Tags:        repo.Topics,
		Stars:       repo.GetStargazersCount(),
		License:     repo.GetLicense().GetName(),
		UpdatedAt:   repo.GetUpdatedAt().Time,
		URL:         repo.GetHTMLURL(),
		Apps:        make([]AppDetails, 0, len(appYML.Apps)),
	}
	for name, info := range appYML.Apps {
		goMod, err := getGoMod(ctx, client, repo, path.Clean(info.Path))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get go.mod for app %s", name)
		}

		result.Apps = append(result.Apps, AppDetails{
			Name:          name,
			Description:   info.Description,
			Path:          info.Path,
			GoVersion:     goMod.Go.Version,
			IgniteVersion: findCLIVersion(goMod),
		})
	}

	return result, nil
}

func validatePackageURL(pkgURL string) (owner, name string, err error) {
	parts := strings.Split(pkgURL, "/")
	if len(parts) != 3 {
		return "", "", fmt.Errorf("package URL must be in github.com/{owner}/{repo} format")
	}
	if parts[0] != "github.com" {
		return "", "", fmt.Errorf("only GitHub packages are supported")
	}

	return parts[1], parts[2], nil
}

func getGoMod(ctx context.Context, client *xgithub.Client, repo *github.Repository, fpath string) (*modfile.File, error) {
	contents, err := client.GetFileContent(ctx, repo.GetOwner().GetLogin(), repo.GetName(), path.Join(fpath, "go.mod"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file content")
	}

	mod, err := modfile.Parse(fmt.Sprintf("%s/%s", repo.GetFullName(), "go.mod"), contents, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse go.mod")
	}

	return mod, nil
}

func findCLIVersion(modFile *modfile.File) string {
	for _, require := range modFile.Require {
		if strings.HasPrefix(require.Mod.Path, igniteCLIPackage) {
			return require.Mod.Version
		}
	}

	return ""
}
